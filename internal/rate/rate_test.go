package rate

import (
	"time"
)

const (
	B  int64 = 1
	KB int64 = B * 1024
	M  int64 = KB * 1024

	deviation      time.Duration = 8
	defaultBufSize               = 4096
)

func getExpectCost(dataSize, bytesPerSecond int64) (expectCost time.Duration, maxCost time.Duration, minCost time.Duration) {
	expectCost = time.Second * time.Duration(dataSize) / time.Duration(bytesPerSecond)

	if bytesPerSecond <= defaultBufSize && dataSize <= defaultBufSize {
		expectCost = time.Second
	} else if bytesPerSecond > defaultBufSize && dataSize <= defaultBufSize {
		expectCost = 0
	}
	maxCost = expectCost * (100 + deviation) / 100
	minCost = expectCost * (100 - deviation) / 100
	// the minimum deviation time is 1 second
	if expectCost-minCost < time.Second {
		minCost = expectCost - time.Second
	}
	if minCost < 0 {
		minCost = 0
	}
	if maxCost <= 0 {
		maxCost = time.Second * deviation / 100
	}
	return expectCost, maxCost, minCost
}

type writer struct {
}

func (w *writer) Write(p []byte) (n int, err error) {
	return len(p), nil
}
