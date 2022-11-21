package mtls

import (
	"fmt"
	"log"
	"os"

	"git.com/balena-labs-research/secure-store/cmd/secure-store/flags"
	"git.com/balena-labs-research/secure-store/pkg/mtls"
)

func GenerateMTLSKeys() {
	fmt.Println("Generating mTLS keys...")

	keyPEM, certPEM, err := mtls.GenerateKeys(4096, 3650, flags.ServerHostname)
	if err != nil {
		if err != nil {
			log.Fatal(err)
		}
	}
	if err := os.WriteFile(flags.KeyPath, keyPEM, 0644); err != nil {
		if err != nil {
			log.Fatal(err)
		}
	}
	if err := os.WriteFile(flags.CertPath, certPEM, 0644); err != nil {
		if err != nil {
			log.Fatal(err)
		}
	}
}

func ValidateMTLSKeys() {
	// If key file doesn't exist, but environment variables do, then generate the key files
	// from the env vars
	_, checkCert := os.Stat(flags.CertPath)
	_, checkKey := os.Stat(flags.KeyPath)

	// TODO: Use of environment variables for storing MTLS keys is untested and subsequently
	// undocumented.
	encryptKey := os.Getenv("SECURE_KEY")
	encryptCert := os.Getenv("SECURE_CERT")

	// If not requesting new keys, and if current key files do not exist already. Files
	// take precedent to avoid attempts to circumvent keys in event of env vars being
	// compromised
	if os.IsNotExist(checkCert) || os.IsNotExist(checkKey) {
		// Check for env variables and generate keys if they exist
		if encryptKey != "" && encryptCert != "" {
			if err := os.WriteFile(flags.CertPath, []byte(encryptCert), 0644); err != nil {
				log.Fatal(err)
			}
			if err := os.WriteFile(flags.KeyPath, []byte(encryptKey), 0644); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatal("There are no MTLS keys or certificates available")
		}
	}
}
