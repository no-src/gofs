package main

import (
	"flag"
	"github.com/no-src/gofs/monitor"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/sync"
	"github.com/no-src/gofs/version"
	"github.com/no-src/log"
	"time"
)

var (
	SrcPath      string
	TargetPath   string
	LogLevel     int
	LogDir       string
	FileLogger   bool
	RetryCount   int
	RetryWait    time.Duration
	BufSize      int
	PrintVersion bool
)

func main() {
	flag.BoolVar(&PrintVersion, "v", false, "print version info")
	flag.StringVar(&SrcPath, "src", "", "source path by monitor")
	flag.StringVar(&TargetPath, "target", "", "target path to backup")
	flag.IntVar(&LogLevel, "log_level", int(log.InfoLevel), "set log level, default is INFO. DEBUG=0 INFO=1 WARN=2 ERROR=3")
	flag.BoolVar(&FileLogger, "file_log", false, "enable file logger")
	flag.StringVar(&LogDir, "log_dir", "./logs/", "set log file's dir")
	flag.IntVar(&RetryCount, "retry_count", 15, "if execute failed, then retry to work retry_count times")
	flag.DurationVar(&RetryWait, "retry_wait", time.Second*5, "if retry to work, wait retry_wait time then do")
	flag.IntVar(&BufSize, "buf_size", 1024*1024, "read and write buffer byte size")
	flag.Parse()

	if PrintVersion {
		version.PrintVersionInfo()
		return
	}

	// init logger
	var loggers []log.Logger
	loggers = append(loggers, log.NewConsoleLogger(log.Level(LogLevel)))
	if FileLogger {
		loggers = append(loggers, log.NewFileLogger(log.Level(LogLevel), LogDir, "gofs"))
	}
	log.InitDefaultLogger(log.NewMultiLogger(loggers...))
	defer log.Log("gofs exited!")

	// create syncer
	syncer, err := sync.NewDiskSync(SrcPath, TargetPath, BufSize)
	if err != nil {
		log.Error(err, "create DiskSync error")
		return
	}

	// create retry
	retry := retry.NewRetry(RetryCount, RetryWait)

	// create monitor
	monitor, err := monitor.NewFsNotifyMonitor(syncer, retry)
	if err != nil {
		log.Error(err, "create fsNotifyMonitor error")
		return
	}
	defer monitor.Close()

	// add to monitor
	err = monitor.Monitor(SrcPath)
	if err != nil {
		log.Error(err, "monitor error, program will be exit")
		return
	}

	// start monitor
	log.Log("file monitor is starting...")
	err = monitor.Start()
	if err != nil {
		log.Log("start to monitor failed, %s", err.Error())
	}
}
