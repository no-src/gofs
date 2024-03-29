package cmd

import (
	"github.com/no-src/gofs/conf"
	"github.com/no-src/gofs/logger"
	"github.com/no-src/log"
	"github.com/no-src/log/formatter"
	"github.com/no-src/log/level"
	"github.com/no-src/log/option"
)

var (
	innerLogger = logger.InnerLogger()
)

const (
	debugLogLevel = level.DebugLevel
)

// initDefaultLogger init the default logger
func initDefaultLogger(c conf.Config) (*logger.Logger, error) {
	// init log formatter
	if c.LogFormat != formatter.TextFormatter {
		innerLogger.Info("switch logger format to %s", c.LogFormat)
	}

	var loggers []log.Logger
	loggers = append(loggers, log.NewConsoleLogger(level.Level(c.LogLevel)))
	if c.EnableFileLogger {
		filePrefix := "gofs_"
		if c.IsDaemon {
			filePrefix += "daemon_"
		}
		flogger, err := log.NewFileLoggerWithOption(option.NewFileLoggerOption(level.Level(c.LogLevel), c.LogDir, filePrefix, c.LogFlush, c.LogFlushInterval.Duration(), c.LogSplitDate))
		if err != nil {
			innerLogger.Error(err, "init file logger error")
			return nil, err
		}
		loggers = append(loggers, flogger)
	}

	defaultLogger := log.NewMultiLogger(loggers...).WithFormatter(formatter.New(c.LogFormat))
	return logger.NewLogger(defaultLogger, log.NewDefaultSampleLogger(defaultLogger, c.LogSampleRate)), nil
}

// initWebServerLogger init the web server logger
func initWebServerLogger(c conf.Config) (*logger.Logger, error) {
	var webLogger = log.NewConsoleLogger(level.Level(c.LogLevel))
	if c.EnableFileLogger && c.EnableFileServer {
		webFileLogger, err := log.NewFileLoggerWithOption(option.NewFileLoggerOption(level.Level(c.LogLevel), c.LogDir, "web_", c.LogFlush, c.LogFlushInterval.Duration(), c.LogSplitDate))
		if err != nil {
			innerLogger.Error(err, "init the web server file logger error")
			return nil, err
		}
		webLogger = log.NewMultiLogger(webFileLogger, webLogger).WithFormatter(formatter.New(c.LogFormat))
	}
	return logger.NewLogger(webLogger, log.NewDefaultSampleLogger(webLogger, c.LogSampleRate)), nil
}

// initEventLogger init the event logger
func initEventLogger(c conf.Config) (log.Logger, error) {
	var eventLogger = log.NewEmptyLogger()
	if c.EnableEventLog {
		eventFileLogger, err := log.NewFileLoggerWithOption(option.NewFileLoggerOption(level.Level(c.LogLevel), c.LogDir, "event_", c.LogFlush, c.LogFlushInterval.Duration(), c.LogSplitDate))
		if err != nil {
			innerLogger.Error(err, "init the event file logger error")
			return nil, err
		}
		eventLogger = eventFileLogger.WithFormatter(formatter.New(c.LogFormat))
	}
	return eventLogger, nil
}
