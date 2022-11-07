package flags

import "flag"

var (
	CertPath       string
	ConfigPath     string
	DecryptEnvOnly bool
	EncryptFolder  string
	EncryptString  string
	ForceUnmount   bool
	GenerateKeys   bool
	GenerateNewKey bool
	KeyPath        string
	MountPoint     string
	Port           string
	ServerHostname string
	StartClient    bool
	UserKey        string
	StartServer    bool
)

func init() {
	// Set flags
	flag.StringVar(&CertPath, "certificate-path", "./cert.pem", "Path for the MTLS certificate")
	flag.StringVar(&EncryptString, "string", "", "String to be encrypted. This flag requires the `-key` flag")
	flag.BoolVar(&GenerateKeys, "generate-keys", false, "Generate the two mTLS keys and save them as files")
	flag.BoolVar(&GenerateNewKey, "new-key", false, "Generate a new key for use in your secure store")
	flag.StringVar(&KeyPath, "key-path", "./key.pem", "Path for the MTLS key")
	flag.StringVar(&ServerHostname, "hostname", "localhost", "Override the default hostname for the server and MTLS key configuration")
	flag.StringVar(&UserKey, "key", "", "A key to encrypt the string passed via `-string`. A random key can be generated using `-new-key`")
}
