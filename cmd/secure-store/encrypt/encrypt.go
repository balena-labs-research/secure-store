package encrypt

import (
	"encoding/hex"
	"fmt"
	"log"

	"git.com/maggie0002/secure-store/pkg/encode"
	"git.com/maggie0002/secure-store/pkg/encryption"
)

func EncryptString(key string, text string) {
	encrypted, err := encryption.EncryptString(text, key)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Encrypted string:")
	fmt.Println("\033[34m", encrypted)
	fmt.Printf("\033[0m")
}

func GenerateNewKey() {
	bytes, err := encode.GenerateRandomBytes(32)

	if err != nil {
		log.Fatal(err)
	}

	key := hex.EncodeToString(bytes)

	fmt.Println("Your new key can be used to encrypt strings using the `-string` flag:")
	fmt.Println("\033[34m", key)
	fmt.Printf("\033[0m")
}