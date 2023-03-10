package rate

import (
	"time"
)

const (
	B  int64 = 1
	KB int64 = B * 1024
	M  int64 = KB * 1024

	deviation      time.Duration = 5
	defaultBufSize               = 4096
)

func getExpectCost(dataSize, bytesPerSecond int64) (expectCost time.Duration, max time.Duration, min time.Duration) {
	expectCost = time.Second * time.Duration(dataSize) / time.Duration(bytesPerSecond)

	if bytesPerSecond <= defaultBufSize && dataSize <= defaultBufSize {
		expectCost = time.Second
	} else if bytesPerSecond > defaultBufSize && dataSize <= defaultBufSize {
		expectCost = 0
	}
	max = expectCost * (100 + deviation) / 100
	min = expectCost * (100 - deviation) / 100
	// the minimum deviation time is 1 second
	if expectCost-min < time.Second {
		min = expectCost - time.Second
	}
	if min < 0 {
		min = 0
	}
	if max <= 0 {
		max = time.Second * deviation / 100
	}
	return expectCost, max, min
}

type writer struct {
}

func (w *writer) Write(p []byte) (n int, err error) {
	return len(p), nil
}
