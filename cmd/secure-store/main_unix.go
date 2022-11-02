//go:build linux

package main

import (
	"flag"
	"fmt"

	"git.com/maggie0002/secure-store/cmd/secure-store/decrypt"
	"git.com/maggie0002/secure-store/cmd/secure-store/encryption"
	"git.com/maggie0002/secure-store/cmd/secure-store/flags"
	"git.com/maggie0002/secure-store/cmd/secure-store/mount"
	"git.com/maggie0002/secure-store/cmd/secure-store/mtls"
	"git.com/maggie0002/secure-store/cmd/secure-store/server"
)

func main() {
	// Parse all flags from all files
	flag.Parse()

	// Take action based on flag
	switch {
	case flags.StartMount:
		mtls.ValidateKeys()
		decrypt.StartMount()
		mount.UnMount(flags.MountPoint)
	case flags.GenerateKeys:
		mtls.GenerateMTLSKeys()
	case flags.StartServer:
		mtls.ValidateKeys()
		server.StartServer()
	case flags.GenerateNewKey:
		encryption.GenerateNewKey()
	case flags.EncryptString != "" && flags.UserKey != "":
		encryption.EncryptString(flags.UserKey, flags.EncryptString)
	case flags.EncryptFolder != "" && flags.UserKey != "":
		mount.EncryptFolder(flags.UserKey, flags.EncryptFolder)
		mount.UnMount(flags.MountPoint)
	default:
		fmt.Println("")
		fmt.Println("Secure Store")
		fmt.Println("---")
		fmt.Println("Pass flags to start the server, mount, or generate keys. " +
			"Your programme to execute after decryption should be passed as arguments (not flags)")
		fmt.Println("---")
		flag.PrintDefaults()
	}
}
