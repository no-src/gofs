package tran

import "errors"

const (
	// LF the string of linefeed
	LF = "\n"
)

var (
	// EndIdentity the network communication end identity
	EndIdentity = []byte("_$#END#$_")
	// ErrorIdentity the network communication error identity
	ErrorIdentity = []byte("_$#ERR#$_")
	// ErrorEndIdentity the network communication end identity with error identity
	ErrorEndIdentity = append(ErrorIdentity, EndIdentity...)
	// LFBytes the linefeed bytes
	LFBytes = []byte(LF)
)

var (
	// ErrServerExecute the remote server execute error
	ErrServerExecute = errors.New("server execute error")
)
