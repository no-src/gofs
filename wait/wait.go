package wait

// WaitDone support execute the work synchronously and mark the work done
type WaitDone interface {
	Wait
	// Done mark the work execute finished
	Done()
	// DoneWithError mark the work execute finished with error info
	DoneWithError(err error)
}

// Wait the interface that implements execute work synchronously
type Wait interface {
	// Wait wait to the work execute finished
	Wait() error
}

// NewWaitDone create an instance of WaitOne to support execute the work synchronously
func NewWaitDone() WaitDone {
	w := &wait{
		c: make(chan error, 1),
	}
	return w
}

type wait struct {
	c chan error
}

func (w *wait) Wait() error {
	return <-w.c
}

func (w *wait) Done() {
	w.c <- nil
}

func (w *wait) DoneWithError(err error) {
	w.c <- err
}
