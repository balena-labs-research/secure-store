package mount

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"git.com/maggie0002/secure-store/cmd/secure-store/flags"
	"git.com/maggie0002/secure-store/pkg/encode"
	"git.com/maggie0002/secure-store/pkg/rclone"
	"github.com/rclone/rclone/fs"
	"github.com/rclone/rclone/fs/config"
	"github.com/rclone/rclone/fs/config/configfile"
	"github.com/rclone/rclone/fs/rc"
)

var (
	mountPoint = flags.MountPoint
)

func CreateMount(passwd string) error {
	// Never ask for passwords, fail instead.
	ci := fs.GetConfig(context.Background())
	ci.AskPassword = false

	// Set config file path and initialize
	err := config.SetConfigPath(flags.ConfigPath)

	if err != nil {
		log.Println(err)
	}

	configfile.Install()

	// Set the config file password
	err = config.SetConfigPassword(passwd)

	if err != nil {
		log.Println(err)
	}

	// Check if .conf exists, and if not then make it
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
		// password which encrypts .conf, and .conf contains this longer password.
		encParams := rc.Params{"remote": "storage", "password": randomPassword}
		encOpts := config.UpdateRemoteOpt{NonInteractive: true, Obscure: true}

		// Create the encrypted storage mount config
		_, err = config.CreateRemote(context.Background(),
			"encrypted_storage",
			"crypt",
			encParams,
			encOpts)

		if err != nil {
			log.Fatal(err.Error())
		}
	}

	// Check if folder exists and make it if not
	_, err = os.Stat(mountPoint)
	if os.IsNotExist(err) {
		err := os.Mkdir(mountPoint, os.ModeDir)

		if err != nil {
			log.Fatal(err.Error())
		}

	}

	// Attempt to remove any existing dangling mount point or socket. Only called
	// when --force-unmount is called
	if flags.ForceUnmount && !os.IsNotExist(err) {
		unMount(mountPoint)
	}

	// Parameters for the rclone mount
	mountParams := rc.Params{
		"fs":         "encrypted_storage:",
		"mountPoint": mountPoint,
		"mountType":  "mount",
		"mountOpt":   "{\"AllowNonEmpty\": true}"} // AllowNonEmpty is required to mount on to an existing directory

	// Mount the volumes
	_, err = rclone.CreateRcloneMount(context.Background(), mountParams)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\033[34m", "Mount created at "+mountPoint)
	fmt.Printf("\033[0m")

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

func unMount(mountPoint string) {
	// Force unmount the directory. Dangling mounts can occur if not exited
	// cleanly leaving FUSE sockets open with nothing listening on the other end.
	err := exec.Command("fusermount", "-u", "-z", mountPoint).Run()
	if err != nil {
		log.Println(err)
		log.Println("Mount point may not exist. Continuing...")
	}
}
