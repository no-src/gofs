package tran

import "errors"

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
