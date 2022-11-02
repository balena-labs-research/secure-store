//go:build !linux

package main

import (
	"flag"
	"fmt"

	"git.com/maggie0002/secure-store/cmd/secure-store/encryption"
	"git.com/maggie0002/secure-store/cmd/secure-store/flags"
	"git.com/maggie0002/secure-store/cmd/secure-store/mtls"
)

func main() {
	// Parse all flags from all files
	flag.Parse()

	// Take action based on flag
	switch {
	case flags.GenerateKeys:
		mtls.GenerateMTLSKeys()
	case flags.GenerateNewKey:
		encryption.GenerateNewKey()
	case flags.EncryptString != "" && flags.UserKey != "":
		encryption.EncryptString(flags.UserKey, flags.EncryptString)
	default:
		fmt.Println("")
		fmt.Println("Secure Store")
		fmt.Println("---")
		flag.PrintDefaults()
	}
}
