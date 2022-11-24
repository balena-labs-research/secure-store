package encrypt

import (
	"encoding/hex"
	"fmt"

	log "github.com/sirupsen/logrus"

	"git.com/balena-labs-research/secure-store/pkg/encode"
	"git.com/balena-labs-research/secure-store/pkg/encryption"
)

func EncryptString(key string, text string) string {
	encrypted, err := encryption.EncryptString(text, key)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Encrypted string:")
	fmt.Println("\033[34m", encrypted)
	fmt.Printf("\033[0m")

	return encrypted
}

func GenerateNewKey() string {
	bytes, err := encode.GenerateRandomBytes(32)

	if err != nil {
		log.Fatal(err)
	}

	key := hex.EncodeToString(bytes)

	fmt.Println("Your new key can be used to encrypt strings using the `-string` flag:")
	fmt.Println("\033[34m", key)
	fmt.Printf("\033[0m")

	return key
}
