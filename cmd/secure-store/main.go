//go:build !linux

// This file runs on non-linux systems. It excludes incompatible commands.

package main

import (
	"flag"
	"fmt"

	"git.com/balena-labs-research/secure-store/cmd/secure-store/encrypt"
	"git.com/balena-labs-research/secure-store/cmd/secure-store/flags"
	"git.com/balena-labs-research/secure-store/cmd/secure-store/mtls"
)

func main() {
	// Parse all flags from all files
	flag.Parse()

	// Take action based on flag
	switch {
	case flags.EncryptString != "" && flags.UserPassword != "":
		encrypt.EncryptString(flags.UserPassword, flags.EncryptString)
	case flags.GenerateKeys:
		mtls.GenerateMTLSKeys()
	case flags.GenerateNewKey:
		encrypt.GenerateNewKey()
	default:
		fmt.Println("")
		fmt.Println("Secure Store")
		fmt.Println("---")
		flag.PrintDefaults()
	}
}
