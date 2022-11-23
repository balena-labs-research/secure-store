package server

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"git.com/balena-labs-research/secure-store/cmd/secure-store/flags"
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

func StartServer() {
	// If the RCLONE default env var is used, accept it and add it as the password
	if os.Getenv("STORE_PASSWORD") == "" && os.Getenv("RCLONE_CONFIG_PASS") != "" {
		os.Setenv("STORE_PASSWORD", os.Getenv("RCLONE_CONFIG_PASS"))
	}

	if os.Getenv("STORE_PASSWORD") == "" {
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
