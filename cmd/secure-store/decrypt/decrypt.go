package decrypt

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"git.com/balena-labs-research/secure-store/pkg/encryption"
)

var out []string

func DecryptEnvs(envs []string, password string) []string {
	for _, env := range envs {
		// If first 10 characters match encrypted prefix
		if len(env) > 9 && env[:10] == "ENCRYPTED_" {
			// Split the key and value
			s := strings.Split(env, "=")

			fmt.Println("\033[34m", "Decrypting environment variable "+s[0])
			fmt.Printf("\033[0m")

			// Decrypt the environment variable
			decryptedValue, err := encryption.DecryptString(s[1], password)

			if err != nil {
				log.Println(err)
			}

			// Trim the ENCRYPTED_ prefix from the value
			trimmedValue := strings.TrimPrefix(s[0], "ENCRYPTED_")

			// Set the decrypted environment variable
			os.Setenv(trimmedValue, decryptedValue)

			out = append(out, trimmedValue+"="+decryptedValue)
		}
	}

	return out
}
