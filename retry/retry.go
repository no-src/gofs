package retry

type Retry interface {
	Do(f func() error, desc string)
}
