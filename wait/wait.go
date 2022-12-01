package wait

import "sync"

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
		c: make(chan error, 1),
	}
	return w
}

type waitDone struct {
	c    chan error
	done bool
	mu   sync.Mutex
}

func (w *waitDone) Wait() error {
	return <-w.c
}

func (w *waitDone) Done() {
	w.DoneWithError(nil)
}

func (w *waitDone) DoneWithError(err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if !w.done {
		w.c <- err
		w.done = true
	}
}
