package main

import (
	"flag"
	"github.com/no-src/gofs/daemon"
	"github.com/no-src/gofs/monitor"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/sync"
	"github.com/no-src/gofs/version"
	"github.com/no-src/log"
	"time"
)

var (
	SrcPath            string
	TargetPath         string
	LogLevel           int
	FileLogger         bool
	LogDir             string
	LogFlush           bool
	LogFlushInterval   time.Duration
	RetryCount         int
	RetryWait          time.Duration
	BufSize            int
	PrintVersion       bool
	SyncOnce           bool
	Daemon             bool
	DaemonPid          bool
	DaemonDelay        time.Duration
	DaemonMonitorDelay time.Duration
	KillPPid           bool
	IsSubprocess       bool
)

func main() {
	flag.BoolVar(&PrintVersion, "v", false, "print version info")
	flag.StringVar(&SrcPath, "src", "", "source path by monitor")
	flag.StringVar(&TargetPath, "target", "", "target path to backup")
	flag.IntVar(&LogLevel, "log_level", int(log.InfoLevel), "set log level, default is INFO. DEBUG=0 INFO=1 WARN=2 ERROR=3")
	flag.BoolVar(&FileLogger, "log_file", false, "enable file logger")
	flag.StringVar(&LogDir, "log_dir", "./logs/", "set log file's dir")
	flag.BoolVar(&LogFlush, "log_flush", false, "enable auto flush log with interval")
	flag.DurationVar(&LogFlushInterval, "log_flush_interval", time.Second*3, "set log flush interval duration, you need to enable log_flush first")
	flag.IntVar(&RetryCount, "retry_count", 15, "if execute failed, then retry to work retry_count times")
	flag.DurationVar(&RetryWait, "retry_wait", time.Second*5, "if retry to work, wait retry_wait time then do")
	flag.IntVar(&BufSize, "buf_size", 1024*1024, "read and write buffer byte size")
	flag.BoolVar(&SyncOnce, "sync_once", false, "sync src directory to target directory once")
	flag.BoolVar(&Daemon, "daemon", false, "enable daemon to create and monitor a subprocess to work, you can use [go build -ldflags=\"-H windowsgui\"] to build on Windows")
	flag.BoolVar(&DaemonPid, "daemon_pid", false, "record parent process pid, daemon process pid and worker process pid to pid file")
	flag.DurationVar(&DaemonDelay, "daemon_delay", time.Second, "daemon work interval, wait to create subprocess")
	flag.DurationVar(&DaemonMonitorDelay, "daemon_monitor_delay", time.Second*3, "daemon monitor work interval, wait to check subprocess state")
	flag.BoolVar(&KillPPid, "kill_ppid", false, "try to kill the parent process when it's running")
	flag.BoolVar(&IsSubprocess, daemon.SubprocessTag, false, "tag current process is subprocess")
	flag.Parse()

	// if current is subprocess, then reset the "kill_ppid" and "daemon"
	if IsSubprocess {
		KillPPid = false
		Daemon = false
	}

	if PrintVersion {
		version.PrintVersionInfo()
		return
	}

	// init logger
	var loggers []log.Logger
	loggers = append(loggers, log.NewConsoleLogger(log.Level(LogLevel)))
	if FileLogger {
		filePrefix := "gofs_"
		if Daemon {
			filePrefix += "daemon_"
		}
		loggers = append(loggers, log.NewFileLoggerWithAutoFlush(log.Level(LogLevel), LogDir, filePrefix, LogFlush, LogFlushInterval))
	}
	log.InitDefaultLogger(log.NewMultiLogger(loggers...))
	defer log.Close()

	// kill parent process
	if KillPPid {
		daemon.KillPPid()
	}

	// start the daemon
	if Daemon {
		daemon.Daemon(DaemonPid, DaemonDelay, DaemonMonitorDelay)
		log.Log("daemon exited")
		return
	}

	// create syncer
	syncer, err := sync.NewDiskSync(SrcPath, TargetPath, BufSize)
	if err != nil {
		log.Error(err, "create DiskSync error")
		return
	}

	// process sync once
	if SyncOnce {
		err = syncer.SyncOnce()
		if err != nil {
			log.Error(err, "sync once error")
		} else {
			log.Log("sync once done!")
		}
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
	defer log.Log("gofs exited!")
	err = monitor.Start()
	if err != nil {
		log.Error(err, "start to monitor failed")
	}
}
