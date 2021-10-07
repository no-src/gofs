package retry

type Retry interface {
	// Do if execute failed, retry execute many times
	Do(f func() error, desc string)
}
