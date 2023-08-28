package core

// SSHConfig the config info for SSH authentication
type SSHConfig struct {
	// Username the username for SSH
	Username string
	// Password the password for SSH
	Password string
	// Key the key file for SSH
	Key string
	// KeyPass the passphrase for the key file
	KeyPass string
	// HostKey the host key file used to validate the server's host key
	HostKey string
}
