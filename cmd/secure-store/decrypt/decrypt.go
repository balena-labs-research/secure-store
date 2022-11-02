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
	"time"

	"git.com/maggie0002/secure-store/cmd/secure-store/flags"
	"git.com/maggie0002/secure-store/cmd/secure-store/mount"
)

const (
	delay = 30
)

type Password struct {
	Password string `json:"password"`
}

func StartMount() {
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
			json.Unmarshal(body, &response)

			// Create the mount
			err = mount.CreateMount(response.Password)

			if err != nil {
				log.Println(err)
			}

			// Decrypt any encrypted environment variables
			err = mount.DecryptEnvs(response.Password)

			if err != nil {
				log.Println(err)
			}

			// Execute the passed Args
			fmt.Println("Executing the passed arguments...")
			err = mount.ExecuteArgs()

			if err != nil {
				log.Println(err)
			}

			return
		}
		// Sleep for 5 seconds
		fmt.Println("Unsuccessful request:", err)
		fmt.Println("Retrying in", delay, "seconds...")
		time.Sleep(delay * time.Second)
	}
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
