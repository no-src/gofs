package wait

import (
	"sync"
)

// WaitDone support to execute the work synchronously and mark the work as done
type WaitDone interface {
	Wait
	Done
}

// Wait the interface that implements execute work synchronously
type Wait interface {
	// Wait wait to the work execute finished
	Wait() error
}

// Done support to mark the work as done
type Done interface {
	// Done mark the work execute finished
	Done()
	// DoneWithError mark the work execute finished with error info
	DoneWithError(err error)
}

// NewWaitDone create an instance of WaitDone to support execute the work synchronously
func NewWaitDone() WaitDone {
	w := &waitDone{
		c: make(chan struct{}, 1),
	}
	return w
}

type waitDone struct {
	c    chan struct{}
	mu   sync.Mutex
	done bool
	err  error
}

func (w *waitDone) Wait() error {
	w.mu.Lock()
	done := w.done
	w.mu.Unlock()
	if done {
		return w.err
	}
	<-w.c
	return w.err
}

func (w *waitDone) Done() {
	w.DoneWithError(nil)
}

func (w *waitDone) DoneWithError(err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if !w.done {
		w.done = true
		w.err = err
		// the channel must be closed at the end
		close(w.c)
	}
}
