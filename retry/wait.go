package retry

// WaitDone support execute the work synchronously
type WaitDone interface {
	Wait
	// Done mark the work execute finished
	Done()
}

type Wait interface {
	// Wait wait to the work execute finished
	Wait()
}

// NewWaitDone create an instance of WaitOne to support execute the work synchronously
func NewWaitDone() WaitDone {
	w := &wait{
		c: make(chan bool, 1),
	}
	return w
}

type wait struct {
	c chan bool
}

func (w *wait) Wait() {
	<-w.c
}

func (w *wait) Done() {
	w.c <- true
}
