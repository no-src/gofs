package tran

import "errors"

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

const (
	// LF the string of linefeed
	LF = "\n"
)

var (
	// ErrServerExecute the remote server execute error
	ErrServerExecute = errors.New("server execute error")
)
