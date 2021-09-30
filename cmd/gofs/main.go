package main

import (
	"github.com/no-src/gofs/daemon"
	"github.com/no-src/gofs/monitor"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/sync"
	"github.com/no-src/gofs/version"
	"github.com/no-src/log"
)

func main() {
	// parse all flags
	parseFlags()

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

	// if enable daemon, start a worker to process the following

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

	// start a file server
	go func() {
		if FileServer {
			log.Log("file server [%s] starting...", FileServerAddr)
			err := server.StartFileServer(SrcPath, TargetPath, FileServerAddr)
			if err != nil {
				log.Error(err, "start file server [%s] error", FileServerAddr)
			}
		}
	}()

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
