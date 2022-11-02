package flags

import "flag"

var (
	GenerateKeys   bool
	GenerateNewKey bool
	UserKey        string
	CertPath       string
	ConfigPath     string
	EncryptFolder  string
	EncryptString  string
	KeyPath        string
	StartMount     bool
	ForceUnmount   bool
	MountPoint     string
	Port           string
	ServerHostname string
	StartServer    bool
)

func init() {
	// Set flags
	flag.BoolVar(&GenerateNewKey, "new-key", false, "Generate a new key for use in your secure store")
	flag.StringVar(&UserKey, "key", "", "A key to encrypt the string passed via `-string`. A random key can be generated using `-new-key`")
	flag.BoolVar(&GenerateKeys, "generate-keys", false, "Generate the two mTLS keys and save them as files")
	flag.StringVar(&EncryptString, "string", "", "String to be encrypted. This flag requires the `-key` flag")
	flag.StringVar(&ServerHostname, "hostname", "localhost", "Override the default hostname for the server and MTLS key configuration")
	flag.StringVar(&CertPath, "certificate-path", "./cert.pem", "Path for the MTLS certificate")
	flag.StringVar(&KeyPath, "key-path", "./key.pem", "Path for the MTLS key")
}
