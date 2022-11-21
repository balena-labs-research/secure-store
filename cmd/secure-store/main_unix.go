//go:build linux

// This file runs on linux only.

package main

import (
	"flag"
	"fmt"

	"git.com/balena-labs-research/secure-store/cmd/secure-store/decrypt"
	"git.com/balena-labs-research/secure-store/cmd/secure-store/encrypt"
	"git.com/balena-labs-research/secure-store/cmd/secure-store/flags"
	"git.com/balena-labs-research/secure-store/cmd/secure-store/mount"
	"git.com/balena-labs-research/secure-store/cmd/secure-store/mtls"
	"git.com/balena-labs-research/secure-store/cmd/secure-store/server"
)

func main() {
	// Parse all flags from all files
	flag.Parse()

	// Take action based on flag
	switch {
	case flags.EncryptFolder != "" && flags.UserKey != "":
		mount.EncryptFolder(flags.UserKey, flags.EncryptFolder)
	case flags.EncryptString != "" && flags.UserKey != "":
		encrypt.EncryptString(flags.UserKey, flags.EncryptString)
	case flags.GenerateKeys:
		mtls.GenerateMTLSKeys()
	case flags.GenerateNewKey:
		encrypt.GenerateNewKey()
	case flags.StartClient:
		mtls.ValidateMTLSKeys()
		decrypt.StartClient()
	case flags.StartServer:
		mtls.ValidateMTLSKeys()
		server.StartServer()
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
