//go:build linux

// This file runs on linux only.

package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"

	"git.com/balena-labs-research/secure-store/cmd/secure-store/encrypt"
	"git.com/balena-labs-research/secure-store/cmd/secure-store/flags"
	"git.com/balena-labs-research/secure-store/cmd/secure-store/http"
	"git.com/balena-labs-research/secure-store/cmd/secure-store/mount"
	"git.com/balena-labs-research/secure-store/cmd/secure-store/mtls"
)

func main() {
	// Parse all flags from all files
	flag.Parse()

	// Take action based on flag
	switch {
	case flags.EncryptFolder != "" && flags.UserPassword != "":
		mount.EncryptFolder(flags.UserPassword, flags.EncryptFolder)
	case flags.EncryptString != "" && flags.UserPassword != "":
		encrypt.EncryptString(flags.UserPassword, flags.EncryptString)
	case flags.GenerateKeys:
		mtls.GenerateMTLSKeys()
	case flags.GenerateNewKey:
		encrypt.GenerateNewKey()
	case flags.LocalMount != "":
		mount.LocalMount(flags.LocalMount)
	case flags.StartClient:
		mtls.ValidateMTLSKeys(os.Getenv("MTLS_CERT"), os.Getenv("MTLS_KEY"))
		err := http.StartClient()

		if err != nil {
			// Raise non-zero exit code to ensure Docker's restart on failure policy works
			log.Fatal(err)
		}
	case flags.StartServer:
		mtls.ValidateMTLSKeys(os.Getenv("MTLS_CERT"), os.Getenv("MTLS_KEY"))
		http.StartServer(os.Getenv("STORE_PASSWORD"), os.Getenv("RCLONE_CONFIG_PASS"))
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
