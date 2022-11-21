package mtls

import (
	"encoding/base64"
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

	// If base64 flag is passed then print keys as base64 and do not write the files
	if flags.Base64 {
		fmt.Println("\033[34m", "MTLS_KEY:")
		fmt.Printf("\033[0m")

		fmt.Println(base64.StdEncoding.EncodeToString([]byte(keyPEM)))

		fmt.Println("\033[34m", "MTLS_CERT:")
		fmt.Printf("\033[0m")

		fmt.Println(base64.StdEncoding.EncodeToString([]byte(certPEM)))
	} else {
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
}

func ValidateMTLSKeys() {
	_, checkCert := os.Stat(flags.CertPath)
	_, checkKey := os.Stat(flags.KeyPath)

	encryptCert := os.Getenv("MTLS_CERT")
	encryptKey := os.Getenv("MTLS_KEY")

	// Check if env variables exist and generate keys. Takes precedent over files existing.
	if encryptCert != "" && encryptKey != "" {
		// Base64 decoding of the cert
		cert, err := base64.StdEncoding.DecodeString(encryptCert)
		if err != nil {
			fmt.Printf("Error decoding cert: %s ", err.Error())
			log.Fatal(err)
		}

		// Base64 decoding of the key
		key, err := base64.StdEncoding.DecodeString(encryptKey)
		if err != nil {
			fmt.Printf("Error decoding key: %s ", err.Error())
			log.Fatal(err)
		}

		if err := os.WriteFile(flags.CertPath, []byte(cert), 0644); err != nil {
			log.Fatal(err)
		}
		if err := os.WriteFile(flags.KeyPath, []byte(key), 0644); err != nil {
			log.Fatal(err)
		}
	} else if os.IsNotExist(checkCert) || os.IsNotExist(checkKey) {
		log.Fatal("There are no mTLS keys or certificates available")
	}
}
