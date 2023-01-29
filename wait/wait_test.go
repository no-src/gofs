package wait

import (
	"errors"
	"sync"
	"testing"
	"time"
)

var (
	errWaitDone = errors.New("wait done error")
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

func TestWaitDone_AlreadyDone(t *testing.T) {
	wd := NewWaitDone()
	wd.Done()
	err := wd.Wait()
	if err != nil {
		t.Errorf("test wait done error =>%s", err)
	}
}

func TestWaitDone_WaitManyTimes(t *testing.T) {
	wd := NewWaitDone()
	wg := sync.WaitGroup{}
	count := 10
	wg.Add(count)
	go func() {
		time.After(time.Millisecond)
		for i := 0; i < count; i++ {
			go func() {
				if err := wd.Wait(); err != nil {
					t.Errorf("test wait done error =>%s", err)
				}
				wg.Done()
			}()
		}
		wd.Done()
	}()
	err := wd.Wait()
	if err != nil {
		t.Errorf("test wait done error =>%s", err)
	}
	wg.Wait()
}

func TestWaitDone_DoneManyTimes(t *testing.T) {
	wd := NewWaitDone()
	wg := sync.WaitGroup{}
	count := 10
	wg.Add(count)
	go func() {
		time.After(time.Millisecond)
		for i := 0; i < count; i++ {
			wd.Done()
			wg.Done()
		}
	}()
	err := wd.Wait()
	if err != nil {
		t.Errorf("test wait done error =>%s", err)
	}
	wg.Wait()
}

func TestWaitDone_WaitAndDoneManyTimes(t *testing.T) {
	wd := NewWaitDone()
	wg := sync.WaitGroup{}
	count := 100
	wg.Add(count*2 + 1)
	go func() {
		time.After(time.Millisecond)
		for i := 0; i < count; i++ {
			go func() {
				if err := wd.Wait(); err != nil {
					t.Errorf("test wait done error =>%s", err)
				}
				wg.Done()
			}()
		}
		wd.Done()
		wg.Done()
		for i := 0; i < count; i++ {
			go func() {
				wd.DoneWithError(errWaitDone)
				wg.Done()
			}()
		}
	}()
	err := wd.Wait()
	if err != nil {
		t.Errorf("test wait done error =>%s", err)
	}
	wg.Wait()
}
