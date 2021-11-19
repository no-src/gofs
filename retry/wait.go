package retry

type WaitFinish interface {
	Wait
	Finish()
}

type Wait interface {
	Wait()
}

func NewWaitFinish() WaitFinish {
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

func (w *wait) Finish() {
	w.c <- true
}
