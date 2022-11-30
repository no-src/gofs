package cmd

import (
	"os"

	"github.com/no-src/gofs/internal/signal"
	"github.com/no-src/gofs/wait"
)

// Result the running result of the program
type Result struct {
	init wait.WaitDone
	wd   wait.WaitDone
	nsc  chan signal.NotifySignal
}

func newResult() Result {
	return Result{
		init: wait.NewWaitDone(),
		wd:   wait.NewWaitDone(),
		nsc:  make(chan signal.NotifySignal, 1),
	}
}

// WaitInit wait for the program to finish initialization
func (r Result) WaitInit() error {
	return r.init.Wait()
}

// Wait wait for the program exit
func (r Result) Wait() error {
	return r.wd.Wait()
}

// Notify send a signal to the program
func (r Result) Notify(s os.Signal) error {
	return (<-r.nsc)(s)
}
