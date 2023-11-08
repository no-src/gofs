package logger

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/no-src/log/level"
)

func TestConsoleLogger(t *testing.T) {
	testLogger("console", NewConsoleLogger(level.DebugLevel, 1.0))
}

func TestTestLogger(t *testing.T) {
	testLogger("test", NewTestLogger())
}

func TestEmptyLogger(t *testing.T) {
	testLogger("empty", NewEmptyLogger())
}

func TestInnerLogger(t *testing.T) {
	testLogger("inner", InnerLogger())
}

func TestInnerLoggerConcurrent(t *testing.T) {
	wg := sync.WaitGroup{}
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(n int) {
			testLogger(fmt.Sprintf("inner[%d]", n), InnerLogger())
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func testLogger(name string, logger *Logger) {
	logger.Debug("%s: hello", name)
	logger.Info("%s: hello", name)
	logger.Warn("%s: hello", name)
	logger.Error(errors.New("test error mock"), "%s: hello", name)
	sample := logger.Sample
	sample.Debug("%s sample: hello", name)
	sample.Info("%s sample: hello", name)
	sample.Warn("%s sample: hello", name)
	sample.Error(errors.New("test error mock"), "%s sample: hello", name)
}
