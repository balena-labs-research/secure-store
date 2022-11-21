//go:build linux

package flags

import "flag"

func init() {
	// Set flags
	flag.StringVar(&ConfigPath, "config-path", "./encrypt.conf", "Path for the config file")
	flag.BoolVar(&DecryptEnvOnly, "env-only", false, "Decrypt the encrypted environment variables, but do not create a mount")
	flag.StringVar(&EncryptFolder, "encrypt-content", "", "Encrypt the content of the provided path. This flag requires the `-key` flag")
	flag.BoolVar(&ForceUnmount, "force-unmount", false, "Attempt a forced unmount of the mount path before mounting")
	flag.StringVar(&LocalMount, "local", "", "Create an encrypted mount locally using the provided password")
	flag.StringVar(&MountPoint, "path", "./decrypted", "Path for your decrypted content")
	flag.StringVar(&Port, "port", "8443", "Override the default port for the client and server")
	flag.BoolVar(&StartClient, "decrypt", false, "Start the mTLS client, look for a decryption "+
		"key on the provided address and port and then attempt a decrypt")
	flag.BoolVar(&StartServer, "server", false, "Start the mTLS server and listen for key requests")
}
