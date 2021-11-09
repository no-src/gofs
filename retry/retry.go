package retry

import "time"

type Retry interface {
	// Do if execute failed, retry execute many times
	Do(f func() error, desc string)
	// RetryCount retry count
	RetryCount() int
	// RetryWait wait time after every retry to fail
	RetryWait() time.Duration
}
