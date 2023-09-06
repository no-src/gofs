package logger

import "github.com/no-src/log"

// Logger an internal logger
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
