package result

import (
	"errors"
	"os"
	"testing"
	"time"
)

var (
	errInit     = errors.New("init error")
	errDone     = errors.New("done error")
	errShutdown = errors.New("shutdown")
)

func TestResult(t *testing.T) {
	r := New()
	go func(r Result) {
		r.InitDone()
		r.Done()
	}(r)

	err := r.WaitInit()
	if err != nil {
		t.Errorf("WaitInit: get an error: %v", err)
		return
	}
	err = r.Wait()
	if err != nil {
		t.Errorf("Wait: get an error: %v", err)
	}
}

func TestResult_Shutdown(t *testing.T) {
	r := New()
	go func(r Result) {
		r.InitDone()
		r.RegisterNotifyHandler(func(s os.Signal, timeout ...time.Duration) error {
			r.DoneWithError(errShutdown)
			return nil
		})
	}(r)

	go func(r Result) {
		if err := r.Shutdown(); err != nil {
			t.Errorf("Shutdown: get an error: %v", err)
		}
	}(r)

	err := r.WaitInit()
	if err != nil {
		t.Errorf("WaitInit: get an error: %v", err)
		return
	}
	err = r.Wait()
	if !errors.Is(err, errShutdown) {
		t.Errorf("Wait: expect to get error: %v, but get %v", errShutdown, err)
	}
}

func TestResult_InitError(t *testing.T) {
	r := New()
	go func(r Result) {
		r.InitDoneWithError(errInit)
	}(r)

	err := r.WaitInit()
	if !errors.Is(err, errInit) {
		t.Errorf("WaitInit: expect to get error: %v, but get %v", errInit, err)
	}
}

func TestResult_DoneError(t *testing.T) {
	r := New()
	go func(r Result) {
		r.InitDone()
		r.DoneWithError(errDone)
	}(r)

	err := r.WaitInit()
	if err != nil {
		t.Errorf("WaitInit: get an error: %v", err)
		return
	}
	err = r.Wait()
	if !errors.Is(err, errDone) {
		t.Errorf("Wait: expect to get error: %v, but get %v", errDone, err)
	}
}
