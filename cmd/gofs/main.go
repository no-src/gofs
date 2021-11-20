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
	parseFlags()

	// if current is subprocess, then reset the "kill_ppid" and "daemon"
	if isSubprocess {
		killPPid = false
		isDaemon = false
	}

	// init logger
	var loggers []log.Logger
	loggers = append(loggers, log.NewConsoleLogger(log.Level(logLevel)))
	if fileLogger {
		filePrefix := "gofs_"
		if isDaemon {
			filePrefix += "daemon_"
		}
		loggers = append(loggers, log.NewFileLoggerWithAutoFlush(log.Level(logLevel), logDir, filePrefix, logFlush, logFlushInterval))
	}
	log.InitDefaultLogger(log.NewMultiLogger(loggers...))
	defer log.Close()

	// print version info
	if printVersion {
		version.PrintVersionInfo()
		return
	}

	// kill parent process
	if killPPid {
		daemon.KillPPid()
	}

	// start the daemon
	if isDaemon {
		daemon.Daemon(daemonPid, daemonDelay, daemonMonitorDelay)
		log.Log("daemon exited")
		return
	}

	// if enable daemon, start a worker to process the following

	// start a file server
	if fileServer {
		waitInit := retry.NewWaitDone()
		go func() {
			err := server.StartFileServer(sourceVFS, targetVFS, fileServerAddr, waitInit)
			if err != nil {
				log.Error(err, "start file server [%s] error", fileServerAddr)
			}
		}()
		waitInit.Wait()
	}

	// create syncer
	syncer, err := sync.NewSync(sourceVFS, targetVFS, bufSize)
	if err != nil {
		log.Error(err, "create DiskSync error")
		return
	}

	// create retry
	retry := retry.NewRetry(retryCount, retryWait, retryAsync)

	// create monitor
	monitor, err := monitor.NewMonitor(syncer, retry, syncOnce)
	if err != nil {
		log.Error(err, "create monitor error")
		return
	}
	defer func() {
		if err = monitor.Close(); err != nil {
			log.Error(err, "close monitor error")
		}
	}()

	// start monitor
	log.Log("file monitor is starting...")
	defer log.Log("gofs exited!")
	defer monitor.Close()
	err = monitor.Start()
	if err != nil {
		log.Error(err, "start to monitor failed")
	}
}
