package rate

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"
	"time"
)

const (
	B  int64 = 1
	KB int64 = B * 1024
	M  int64 = KB * 1024

	deviation      time.Duration = 5
	defaultBufSize               = 4096
)

func TestReader(t *testing.T) {
	testCases := []struct {
		name           string
		dataSize       int64
		bytesPerSecond int64
	}{
		{"bufio reader 10M 10M/sec", M * 10, M * 10},
		{"bufio reader 10M 8M/sec", M * 10, M * 8},
		{"bufio reader 10M 5M/sec", M * 10, M * 5},
		{"bufio reader 10M 3M/sec", M * 10, M * 3},
		{"bufio reader <bytesPerSecond less than defaultBufSize> 1KB 250bytes/sec", KB, B * 250},
		{"bufio reader <bytesPerSecond less than defaultBufSize> 1KB 500bytes/sec", KB, B * 500},
		{"bufio reader <bytesPerSecond less than defaultBufSize> 4KB 1KB/sec", KB * 4, KB},
		{"bufio reader <bytesPerSecond less less than defaultBufSize> 4096 bytes 1KB/sec", B * defaultBufSize, KB},
		{"bufio reader <dataSize less than defaultBufSize> 4096 bytes 5KB/sec", B * defaultBufSize, KB * 5},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			buf.WriteString(strings.Repeat("a", int(tc.dataSize)))
			r := NewReader(buf, tc.bytesPerSecond)
			start := time.Now()
			br := bufio.NewReaderSize(r, defaultBufSize)
			_, err := br.WriteTo(&writer{})
			end := time.Now()
			if err != nil {
				t.Errorf("read data error, %v", err)
				return
			}

			expectCost, max, min := getExpectCost(tc.dataSize, tc.bytesPerSecond)
			actualCost := end.Sub(start)
			if actualCost >= min && actualCost <= max {
				rate := 0.0
				if actualCost > 0 {
					rate = float64(tc.dataSize) / actualCost.Seconds()
				}
				t.Logf("[%s] dataSize=%d bytesPerSecond=%d cost=%s, rate=%.2f bytes/sec", tc.name, tc.dataSize, tc.bytesPerSecond, actualCost, rate)
			} else {
				t.Errorf("[%s] dataSize=%d bytesPerSecond=%d expect cost=%s min=%s max=%s, but actual cost=%s", tc.name, tc.dataSize, tc.bytesPerSecond, expectCost, min, max, actualCost)
			}
		})
	}
}

func TestReadAt(t *testing.T) {
	testCases := []struct {
		name           string
		dataSize       int64
		bytesPerSecond int64
	}{
		{"file reader at 10M 10M/sec", M * 10, M * 10},
		{"file reader at 10M 8M/sec", M * 10, M * 8},
		{"file reader at 10M 5M/sec", M * 10, M * 5},
		{"file reader at 10M 3M/sec", M * 10, M * 3},
		{"file reader at <bytesPerSecond less than defaultBufSize> 1KB 250bytes/sec", KB, B * 250},
		{"file reader at <bytesPerSecond less than defaultBufSize> 1KB 500bytes/sec", KB, B * 500},
		{"file reader at <bytesPerSecond less than defaultBufSize> 4KB 1KB/sec", KB * 4, KB},
		{"file reader at <bytesPerSecond less less than defaultBufSize> 4096 bytes 1KB/sec", B * defaultBufSize, KB},
		{"file reader at <dataSize less than defaultBufSize> 4096 bytes 5KB/sec", B * defaultBufSize, KB * 5},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open("./rate_test.go")
			if err != nil {
				t.Errorf("open file error, %v", err)
				return
			}
			defer f.Close()
			r := NewReaderAt(f, tc.bytesPerSecond)
			start := time.Now()
			var total int64 = 0
			for {
				p := make([]byte, defaultBufSize)
				n, err := r.ReadAt(p, 1)
				if err != nil {
					t.Errorf("read file error, %v", err)
					return
				}
				total += int64(n)
				if total >= tc.dataSize {
					break
				}
			}
			end := time.Now()
			if err != nil {
				t.Errorf("read data error, %v", err)
				return
			}

			expectCost, max, min := getExpectCost(tc.dataSize, tc.bytesPerSecond)
			actualCost := end.Sub(start)
			if actualCost >= min && actualCost <= max {
				rate := 0.0
				if actualCost > 0 {
					rate = float64(tc.dataSize) / actualCost.Seconds()
				}
				t.Logf("[%s] dataSize=%d bytesPerSecond=%d cost=%s, rate=%.2f bytes/sec", tc.name, tc.dataSize, tc.bytesPerSecond, actualCost, rate)
			} else {
				t.Errorf("[%s] dataSize=%d bytesPerSecond=%d expect cost=%s min=%s max=%s, but actual cost=%s", tc.name, tc.dataSize, tc.bytesPerSecond, expectCost, min, max, actualCost)
			}
		})
	}
}

func TestNewReader_DisableOrEnableRate(t *testing.T) {
	testCases := []struct {
		name           string
		bytesPerSecond int64
		expectRate     bool
	}{
		{"disable rate by zero rate", 0, true},
		{"disable rate by negative rate", -1, true},
		{"enable rate", 1, false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			d := NewReader(buf, tc.bytesPerSecond)
			switch d.(type) {
			case *bytes.Buffer:
				if !tc.expectRate {
					t.Errorf("expect to get *reader type, actual get *bytes.Buffer")
				}
			case *reader:
				if tc.expectRate {
					t.Errorf("expect to get *bytes.Buffer type, actual get *reader")
				}
			default:
				t.Errorf("unexpected type")
			}
		})
	}
}

func TestNewReaderAt_DisableOrEnableRate(t *testing.T) {
	testCases := []struct {
		name           string
		bytesPerSecond int64
		expectRate     bool
	}{
		{"disable rate by zero rate", 0, true},
		{"disable rate by negative rate", -1, true},
		{"enable rate", 1, false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := NewReaderAt(os.Stdout, tc.bytesPerSecond)
			switch d.(type) {
			case *os.File:
				if !tc.expectRate {
					t.Errorf("expect to get *reader type, actual get *os.File")
				}
			case *reader:
				if tc.expectRate {
					t.Errorf("expect to get *os.File type, actual get *reader")
				}
			default:
				t.Errorf("unexpected type")
			}
		})
	}
}

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
	return expectCost, max, min
}

type writer struct {
}

func (w *writer) Write(p []byte) (n int, err error) {
	return len(p), nil
}
