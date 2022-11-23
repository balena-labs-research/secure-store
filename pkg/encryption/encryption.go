package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
)

func EncryptString(stringToEncrypt string, key string) (string, error) {
	// Use sha256 hash to ensure correct key length
	hsha2 := sha256.Sum256([]byte(key))

	// Create a new Cipher Block from the key
	block, err := aes.NewCipher(hsha2[:])
	if err != nil {
		return "", err
	}

	// Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the data using aesGCM.Seal
	// Since we don't want to save the nonce somewhere else in this case, we add it as a prefix
	// to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(stringToEncrypt), nil)
	return hex.EncodeToString(ciphertext), nil
}

func DecryptString(encryptedString string, key string) (string, error) {
	enc, _ := hex.DecodeString(encryptedString)
	hsha2 := sha256.Sum256([]byte(key))

	// Create a new Cipher Block from the key
	block, err := aes.NewCipher(hsha2[:])
	if err != nil {
		return "", err
	}

	// Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Get the nonce size
	nonceSize := aesGCM.NonceSize()

	// Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	// Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
