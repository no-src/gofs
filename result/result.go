package result

import (
	"os"
	"syscall"

	"github.com/no-src/gofs/internal/signal"
	"github.com/no-src/gofs/wait"
)

// Result control the running result of the program
type Result interface {
	// WaitInit wait for the program to finish initialization work
	WaitInit() error
	// Wait wait for the program exit
	Wait() error
	// Notify send a signal to the program
	Notify(s os.Signal) error
	// Shutdown shutdown the program
	Shutdown() error
	// InitDone mark the initialization work as finished
	InitDone()
	// InitDoneWithError mark the initialization work as failed
	InitDoneWithError(err error)
	// Done mark the work as finished
	Done()
	// DoneWithError mark the work as failed
	DoneWithError(err error)
	// RegisterNotifyHandler register a signal handler
	RegisterNotifyHandler(ns signal.NotifySignal)
}

type result struct {
	init wait.WaitDone
	wd   wait.WaitDone
	nsc  chan signal.NotifySignal
}

// New create an instance of the Result
func New() Result {
	return &result{
		init: wait.NewWaitDone(),
		wd:   wait.NewWaitDone(),
		nsc:  make(chan signal.NotifySignal, 1),
	}
}

func (r *result) WaitInit() error {
	return r.init.Wait()
}

func (r *result) Wait() error {
	return r.wd.Wait()
}

func (r *result) Notify(s os.Signal) error {
	return (<-r.nsc)(s)
}

func (r *result) Shutdown() error {
	return r.Notify(syscall.SIGQUIT)
}

func (r *result) InitDone() {
	r.init.Done()
}

func (r *result) InitDoneWithError(err error) {
	r.init.DoneWithError(err)
}

func (r *result) Done() {
	r.wd.Done()
}

func (r *result) DoneWithError(err error) {
	r.wd.DoneWithError(err)
}

func (r *result) RegisterNotifyHandler(ns signal.NotifySignal) {
	r.nsc <- ns
}
