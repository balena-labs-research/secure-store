package mount

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"git.com/maggie0002/secure-store/cmd/secure-store/flags"
	"git.com/maggie0002/secure-store/pkg/encode"
	"git.com/maggie0002/secure-store/pkg/encryption"
	_ "github.com/rclone/rclone/backend/crypt"
	_ "github.com/rclone/rclone/backend/local"
	_ "github.com/rclone/rclone/cmd/mount"
	"github.com/rclone/rclone/cmd/mountlib"
	"github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/fs/config"
	"github.com/rclone/rclone/fs/config/configfile"
	"github.com/rclone/rclone/fs/rc"
	"github.com/rclone/rclone/vfs/vfsflags"
)

var (
	liveMounts = map[string]*mountlib.MountPoint{}
	mountMu    sync.Mutex
	mountPoint = flags.MountPoint
)

func CreateMount(passwd string) error {
	// Never ask for passwords, fail instead.
	ci := fs.GetConfig(context.Background())
	ci.AskPassword = false

	// Set config file path and initialize
	config.SetConfigPath(flags.ConfigPath)
	configfile.Install()

	// Set the config file password
	config.SetConfigPassword(passwd)

	// Check if rclone.conf exists, and if not then make it
	if _, err := os.Stat(flags.ConfigPath); os.IsNotExist(err) {
		// Create an encrypted storage configuration file
		fmt.Println("Creating encryption configuration file...")

		// Set parameters for local storage
		localParams := rc.Params{}
		var localOpts = config.UpdateRemoteOpt{NonInteractive: true}

		// Create the local storage mount config
		_, err := config.CreateRemote(context.Background(), "storage", "local", localParams, localOpts)

		if err != nil {
			log.Fatal(err.Error())
		}

		// Generate a random encryption key
		randomPassword, err := encode.GenerateRandomString(32)

		if err != nil {
			log.Fatal(err.Error())
		}

		// Set parameters for encrypted storage. This includes a random password for encrypting
		// the files. This isn't the password required to decrypt, instead that is done by another
		// password which encrypts rclone.conf, and rclone.conf contains this longer password.
		encParams := rc.Params{"remote": "storage", "password": randomPassword}
		encOpts := config.UpdateRemoteOpt{NonInteractive: true, Obscure: true}

		// Create the encrypted storage mount config
		_, err = config.CreateRemote(context.Background(),
			"encrypted_storage",
			"crypt",
			encParams, encOpts)

		if err != nil {
			log.Fatal(err.Error())
		}
	}

	// Parameters for the rclone mount
	mountParams := rc.Params{
		"fs":         "encrypted_storage:",
		"mountPoint": mountPoint,
		"mountType":  "mount",
		"mountOpt":   "{\"AllowNonEmpty\": true}"} // AllowNonEmpty is required to mount on to an existing directory

	// Check if folder exists and make it if not
	_, err := os.Stat(mountPoint)
	if os.IsNotExist(err) {
		os.Mkdir(mountPoint, os.ModeDir)
	}

	// Attempt to remove any existing dangling mount point or socket. Only called
	// when --force-unmount is called
	if flags.ForceUnmount && !os.IsNotExist(err) {
		UnMount(mountPoint)
	}

	// Mount the volumes
	_, err = createRcloneMount(context.Background(), mountParams)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func createRcloneMount(ctx context.Context, in rc.Params) (out rc.Params, err error) {
	mountPoint, err := in.GetString("mountPoint")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	vfsOpt := vfsflags.Opt
	err = in.GetStructMissingOK("vfsOpt", &vfsOpt)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	mountOpt := mountlib.Options{AllowOther: true}
	err = in.GetStructMissingOK("mountOpt", &mountOpt)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	mountType, err := in.GetString("mountType")

	mountMu.Lock()
	defer mountMu.Unlock()

	if err != nil {
		mountType = ""
	}
	_, mountFn := mountlib.ResolveMountMethod(mountType)
	if mountFn == nil {
		return nil, errors.New("mount Option specified is not registered, or is invalid")
	}

	// Get Fs.fs to be mounted from fs parameter in the params
	fdst, err := rc.GetFs(ctx, in)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	mnt := mountlib.NewMountPoint(mountFn, mountPoint, fdst, &mountOpt, &vfsOpt)

	_, err = mnt.Mount()

	if err != nil {
		log.Printf("mount FAILED: %v", err)
		return nil, err
	}

	// Add mount to list if mount point was successfully created
	liveMounts[mountPoint] = mnt

	fmt.Println("\033[34m", "Mount created at "+mountPoint)
	fmt.Printf("\033[0m")
	return nil, nil
}

func DecryptEnvs(password string) error {
	for _, env := range os.Environ() {
		// If first 10 characters match encrypted prefix
		if len(env) > 9 && env[:10] == "ENCRYPTED_" {
			// Split the key and value
			s := strings.Split(env, "=")

			fmt.Println("\033[34m", "Decrypting environment variable "+s[0])
			fmt.Printf("\033[0m")

			// Decrypt the environment variable
			decryptedValue, err := encryption.DecryptString(s[1], password)

			if err != nil {
				log.Println(err)
			}

			// Trim the ENCRYPTED_ prefix from the value
			trimedValue := strings.TrimPrefix(s[0], "ENCRYPTED_")

			// Set the decrypted environment variable
			os.Setenv(trimedValue, decryptedValue)
		}
	}
	return nil
}

func EncryptFolder(key string, sourceFolder string) {
	// Check if source folder exists
	_, err := os.Stat(sourceFolder)
	if os.IsNotExist(err) {
		log.Fatalln("Source folder does not exist")
	}

	err = CreateMount(key)

	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("cp", "-r", sourceFolder, flags.MountPoint)
	err = cmd.Run()

	if err != nil {
		log.Printf("Error copying encrypted files: %v", err)
	}
}

func ExecuteArgs() error {
	// Execute the main process requested by the user. It is run here as this app needs to
	// keep running to serve the encryption, and to allow both processes to remain as PID 1

	if len(flag.Args()) == 0 {
		fmt.Println("No arguments were passed to execute.")
		return nil
	}

	cmd := exec.Command(flag.Args()[0], flag.Args()[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		return err
	}

	return nil
}

func UnMount(mountPoint string) {
	// Check if the chosen path is already a mount point. If it is, we assume that it
	// is a dangling mount point and try to unmount and replace it with this new one.
	// Dangling mounts can occur if not exited cleanly, leaving FUSE sockets open
	// with nothing listening on the other end.
	err := exec.Command("mountpoint", "-q", mountPoint).Run()

	if err != nil {
		err := exec.Command("fusermount", "-u", mountPoint).Run()
		if err != nil {
			log.Printf("fusermount failed: %v. Continuing...", err)
		}
	}
}
