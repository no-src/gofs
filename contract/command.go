package contract

type Command []byte

var (
	// InfoCommand send the info command to get server info
	InfoCommand Command = []byte("info")
	// AuthCommand send the auth command with username and password to sign in the server
	AuthCommand Command = []byte("auth")
	// UnknownCommand unknown command
	UnknownCommand Command = []byte("unknown")
)
