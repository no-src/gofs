package logger

import (
	"github.com/no-src/log"
	"github.com/no-src/log/formatter"
	"github.com/no-src/log/level"
)

var (
	// DefaultLevel default log level
	DefaultLevel = level.InfoLevel
	// DefaultFormatter default log formatter
	DefaultFormatter = formatter.TextFormatter
)

// Logger an logger component
type Logger struct {
	// Logger the default logger
	log.Logger

	// Sample the sample logger
	Sample log.Logger
}

// NewLogger create an instance of Logger
func NewLogger(logger, sample log.Logger) *Logger {
	return &Logger{
		Logger: logger,
		Sample: sample,
	}
}

// NewConsoleLogger return a console logger
func NewConsoleLogger(lvl level.Level, sampleRate float64) *Logger {
	logger := log.NewConsoleLogger(lvl)
	sample := log.NewDefaultSampleLogger(logger, sampleRate)
	return NewLogger(logger, sample)
}

// NewTestLogger return a logger used for the test
func NewTestLogger() *Logger {
	return NewConsoleLogger(level.DebugLevel, 1)
}

// NewEmptyLogger get an empty logger, there is nothing to do
func NewEmptyLogger() *Logger {
	logger := log.NewEmptyLogger()
	return NewLogger(logger, logger)
}
