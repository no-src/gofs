package retry

type WaitDone interface {
	Wait
	Done()
}

type Wait interface {
	Wait()
}

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
