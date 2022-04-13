package retry

import (
	"context"
	"time"

	"github.com/no-src/gofs/wait"
)

// Retry if execute return error, then retry to execute with the specified rule
type Retry interface {
	// Do execute work first, if execute failed, retry execute many times
	Do(f func() error, desc string) wait.Wait
	// DoWithContext execute work first, if execute failed, retry execute many times, cancel retry with context
	DoWithContext(ctx context.Context, f func() error, desc string) wait.Wait
	// Count the retry count
	Count() int
	// WaitTime the wait time after every retry to fail
	WaitTime() time.Duration
}
