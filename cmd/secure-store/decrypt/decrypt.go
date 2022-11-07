package decrypt

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"git.com/maggie0002/secure-store/cmd/secure-store/flags"
	"git.com/maggie0002/secure-store/cmd/secure-store/mount"
	"git.com/maggie0002/secure-store/cmd/secure-store/run"
	"git.com/maggie0002/secure-store/pkg/encryption"
)

const (
	delay = 30
)

type Password struct {
	Password string `json:"password"`
}

func DecryptEnvs(password string) error {
	for _, env := range os.Environ() {
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
		}
	}
	return nil
}

func makeRequest(client *http.Client) ([]byte, error) {
	// Request via the HTTPS client over port X
	r, err := client.Get("https://" + flags.ServerHostname + ":" + flags.Port + "/key")
	if err != nil {
		return nil, err
	}

	// Read the response body
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	// Return the body
	return body, nil
}

func StartClient() {
	// Read the key pair to create certificate
	cert, err := tls.LoadX509KeyPair(flags.CertPath, flags.KeyPath)
	if err != nil {
		log.Fatal(err)
	}

	// Create a CA certificate pool and add cert.pem to it
	caCert, err := os.ReadFile(flags.CertPath)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create a HTTPS client and supply the created CA pool and certificate
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
	}

	// Retry the request until it succeeds
	for {
		body, err := makeRequest(client)
		if err == nil {
			fmt.Println("Attempting decrypt...")

			// Decode the JSON
			var response Password
			err = json.Unmarshal(body, &response)

			if err != nil {
				log.Fatal(err)
			}

			// Decrypt any encrypted environment variables
			err = DecryptEnvs(response.Password)

			if err != nil {
				log.Println(err)
			}

			// Create the mount
			if !flags.DecryptEnvOnly {
				err = mount.CreateMount(response.Password)

				if err != nil {
					log.Println(err)
				}
			}

			// Execute the passed Args
			fmt.Println("Executing the passed arguments...")
			err = run.ExecuteArgs()

			if err != nil {
				// Raise non-zero exit code to ensure Docker's restart on failure policy works
				log.Fatal(err)
			}

			return
		}
		// Sleep for 5 seconds
		fmt.Println("Unsuccessful request:", err)
		fmt.Println("Retrying in", delay, "seconds...")
		time.Sleep(delay * time.Second)
	}
}
