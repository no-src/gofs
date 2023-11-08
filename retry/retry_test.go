package retry

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/no-src/gofs/logger"
)

const (
	testRetryCount = 3
	testWaitTime   = time.Millisecond
)

func TestDefaultRetry_SuccessAsync(t *testing.T) {
	testDefaultRetrySuccess(t, true)
}

func TestDefaultRetry_SuccessSync(t *testing.T) {
	testDefaultRetrySuccess(t, false)
}

func testDefaultRetrySuccess(t *testing.T, async bool) {
	tw := testWork{}
	desc := "[async] [do work success] "
	if !async {
		desc = "[sync] [do work success] "
	}
	r := New(testRetryCount, testWaitTime, async, logger.NewTestLogger())
	w := r.Do(tw.doSuccess, desc)
	err := w.Wait()
	if err != nil {
		t.Errorf("%s test retry error, wait return error => %s", desc, err)
	}
	if r.Count() != testRetryCount {
		t.Errorf("%s test retry error, retry count expect:%d, actual:%d", desc, testRetryCount, r.Count())
	}
	if r.WaitTime() != testWaitTime {
		t.Errorf("%s test retry error, retry wait time expect:%s, actual:%s", desc, testWaitTime.String(), r.WaitTime().String())
	}
}

func TestDefaultRetry_FailedAsync(t *testing.T) {
	testDefaultRetryFailed(t, true)
}

func TestDefaultRetry_FailedSync(t *testing.T) {
	testDefaultRetryFailed(t, false)
}

func testDefaultRetryFailed(t *testing.T, async bool) {
	tw := testWork{}
	desc := "[async] [do work failed] "
	if !async {
		desc = "[sync] [do work failed] "
	}
	r := New(testRetryCount, testWaitTime, async, logger.NewTestLogger())
	w := r.Do(tw.doFail, desc)
	err := w.Wait()
	if err != nil {
		t.Errorf("%s test retry error, wait return error => %s", desc, err)
	}
}

func TestDefaultRetry_AbortAsync(t *testing.T) {
	testDefaultRetryAbort(t, true)
}

func TestDefaultRetry_AbortSync(t *testing.T) {
	testDefaultRetryAbort(t, false)
}

func testDefaultRetryAbort(t *testing.T, async bool) {
	tw := testWork{}
	desc := "[async] [do work abort] "
	if !async {
		desc = "[sync] [do work abort] "
	}
	r := New(testRetryCount, testWaitTime, async, logger.NewTestLogger())
	ctx, cancel := context.WithCancel(context.Background())
	w := r.DoWithContext(ctx, tw.doFail, desc)
	cancel()
	err := w.Wait()
	if err != nil {
		t.Errorf("%s test retry error, wait return error => %s", desc, err)
	}
}

func TestDefaultRetry_FailedThenSuccessAsync(t *testing.T) {
	testDefaultRetryFailedThenSuccess(t, true)
}

func TestDefaultRetry_FailedThenSuccessSync(t *testing.T) {
	testDefaultRetryFailedThenSuccess(t, false)
}

func testDefaultRetryFailedThenSuccess(t *testing.T, async bool) {
	tw := testWork{}
	asyncDesc := "[async] "
	if !async {
		asyncDesc = "[sync] "
	}
	r := New(testRetryCount, testWaitTime, async, logger.NewTestLogger())
	w := r.Do(tw.doFailThenSuccess, asyncDesc+"do work failed then success")
	err := w.Wait()
	if err != nil {
		t.Errorf("%s test retry error, wait return error => %s", asyncDesc, err)
	}
}

func TestDefaultRetry_Panic(t *testing.T) {
	tw := testWork{}
	r := New(testRetryCount, testWaitTime, true, logger.NewTestLogger())
	w := r.Do(tw.doPanic, "do work panic")
	err := w.Wait()
	if err != nil {
		t.Errorf("test retry error, wait return error => %s", err)
	}
}

type testWork struct {
	c int
}

func (tw *testWork) doSuccess() error {
	fmt.Println("work done")
	return nil
}

func (tw *testWork) doFail() error {
	return errors.New("work error")
}

func (tw *testWork) doFailThenSuccess() error {
	tw.c++
	if tw.c > 2 {
		return nil
	}
	return errors.New("work error")
}

func (tw *testWork) doPanic() error {
	panic("work panic")
}
