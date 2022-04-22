package wait

import (
	"errors"
	"testing"
	"time"
)

func TestWaitDone(t *testing.T) {
	wd := NewWaitDone()
	go func() {
		time.After(time.Millisecond)
		wd.Done()
	}()
	err := wd.Wait()
	if err != nil {
		t.Errorf("test wait done error =>%s", err)
	}
}

func TestWaitDone_ReturnError(t *testing.T) {
	errWaitDone := errors.New("wait done error")
	wd := NewWaitDone()
	go func() {
		time.After(time.Millisecond)
		wd.DoneWithError(errWaitDone)
	}()
	err := wd.Wait()
	if err == nil || !errors.Is(err, errWaitDone) {
		t.Errorf("test wait done with error failed, expect:%s actual:%s", errWaitDone, err)
	}
}
