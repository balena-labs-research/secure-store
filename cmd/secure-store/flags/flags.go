package flags

import "flag"

var (
	Base64         bool
	CertPath       string
	ConfigPath     string
	DecryptEnvOnly bool
	EncryptFolder  string
	EncryptString  string
	ForceUnmount   bool
	GenerateKeys   bool
	GenerateNewKey bool
	KeyPath        string
	LocalMount     string
	MountPoint     string
	Port           string
	ServerHostname string
	StartClient    bool
	UserPassword   string
	StartServer    bool
)

func init() {
	// Set flags
	flag.BoolVar(&Base64, "base64", false, "Generate base64 outputs instead of files")
	flag.StringVar(&CertPath, "certificate-path", "./cert.pem", "Path for the MTLS certificate")
	flag.StringVar(&EncryptString, "string", "", "String to be encrypted. This flag requires the `-password` flag")
	flag.BoolVar(&GenerateKeys, "generate-keys", false, "Generate the two mTLS keys and save them as files")
	flag.BoolVar(&GenerateNewKey, "new-key", false, "Generate a random key for use in your Secure Store")
	flag.StringVar(&KeyPath, "key-path", "./key.pem", "Path for the MTLS key")
	flag.StringVar(&ServerHostname, "hostname", "localhost", "Override the default hostname for the server and MTLS key configuration")
	flag.StringVar(&UserPassword, "password", "", "The password to encrypt the string passed via `-string`")
}
