package http

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"git.com/balena-labs-research/secure-store/cmd/secure-store/decrypt"
	"git.com/balena-labs-research/secure-store/cmd/secure-store/flags"
	"git.com/balena-labs-research/secure-store/cmd/secure-store/mount"
	"git.com/balena-labs-research/secure-store/cmd/secure-store/run"
)

const (
	delay = 30
)

type Password struct {
	Password string `json:"password"`
}

func keyHandler(w http.ResponseWriter, r *http.Request) {
	response := Password{os.Getenv("STORE_PASSWORD")}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		log.Println(err)
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

func StartClient() error {
	// Read the key pair to create certificate
	cert, err := tls.LoadX509KeyPair(flags.CertPath, flags.KeyPath)
	if err != nil {
		return err
	}

	// Create a CA certificate pool and add cert.pem to it
	caCert, err := os.ReadFile(flags.CertPath)
	if err != nil {
		return err
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
				return err
			}

			// Do not execute the decrypt and mount when running Go tests
			if flag.Lookup("test.v") != nil {
				return nil
			}

			// Decrypt any encrypted environment variables
			decrypt.DecryptEnvs(os.Environ(), response.Password)

			// Create the mount
			if !flags.DecryptEnvOnly {
				err = mount.CreateMount(response.Password)

				if err != nil {
					log.Println(err)
				}
			}

			// Execute the passed Args
			err = run.ExecuteArgs(flag.Args())

			if err != nil {
				// Raise non-zero exit code to ensure Docker's restart on failure policy works
				return err
			}

		}
		// Sleep for 5 seconds
		fmt.Println("Unsuccessful request:", err)
		fmt.Println("Retrying in", delay, "seconds...")
		time.Sleep(delay * time.Second)
	}
}

func StartServer(storePassword string, rcloneConfigPass string) {
	// If the RCLONE default env var is used, accept it and add it as the password
	if storePassword == "" && rcloneConfigPass != "" {
		storePassword = rcloneConfigPass
	}

	if storePassword == "" {
		log.Fatal("No password set in the STORE_PASSWORD environment variable")
	}

	fmt.Println("Server listening for requests...")
	// Set up a resource handler
	http.HandleFunc("/key", keyHandler)

	// Create a CA certificate pool and add cert.pem to it
	caCert, err := os.ReadFile(flags.CertPath)
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create the TLS Config with the CA pool and enable Client certificate validation
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	// Create a Server instance to listen on port 8443 with the TLS config
	server := &http.Server{
		Addr:      ":" + flags.Port,
		TLSConfig: tlsConfig,
	}

	// Listen for HTTPS connections with the server certificate and wait
	log.Fatal(server.ListenAndServeTLS(flags.CertPath, flags.KeyPath))
}
