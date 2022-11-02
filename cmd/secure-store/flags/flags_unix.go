//go:build linux

package flags

import "flag"

func init() {
	// Set flags
	flag.BoolVar(&StartMount, "decrypt", false, "Start the mTLS client, look for a decryption "+
		"key on the provided address and port and then attempt a decrypt")
	flag.StringVar(&Port, "port", "8443", "Override the default port for the client and server")
	flag.BoolVar(&StartServer, "server", false, "Start the mTLS server and listen for key requests")
	flag.BoolVar(&ForceUnmount, "force-unmount", false, "If the mount point already exists, attempt an unmount first")
	flag.StringVar(&MountPoint, "path", "./decrypted", "Path for your decrypted content")
	flag.StringVar(&ConfigPath, "config-path", "./encrypt.conf", "Path for the config file")
	flag.StringVar(&EncryptFolder, "encrypt-content", "", "Encrypt the content of the provided path. This flag requires the `-key` flag")
}
