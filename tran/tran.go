package tran

import "errors"

const (
	DefaultServerHost = "127.0.0.1"
	DefaultServerPort = 9016
)

var (
	EndIdentity      = []byte("_$#END#$_")
	ErrorIdentity    = []byte("_$#ERR#$_")
	ErrorEndIdentity = append(ErrorIdentity, EndIdentity...)
	LFBytes          = []byte(LF)
)

const (
	LF = "\n"
)

var (
	ServerExecuteError = errors.New("server execute error")
)
