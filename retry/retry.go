package retry

import (
	"context"
	"time"
)

// Retry if execute return error, then retry to execute with the specified rule
type Retry interface {
	// Do execute work first, if execute failed, retry execute many times
	Do(f func() error, desc string) Wait
	// DoWithContext execute work first, if execute failed, retry execute many times, cancel retry with context
	DoWithContext(ctx context.Context, f func() error, desc string) Wait
	// RetryCount retry count
	RetryCount() int
	// RetryWait wait time after every retry to fail
	RetryWait() time.Duration
}
