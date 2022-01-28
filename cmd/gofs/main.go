package main

import (
	"github.com/no-src/gofs/about"
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/daemon"
	"github.com/no-src/gofs/internal/signal"
	"github.com/no-src/gofs/monitor"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/server/httpfs"
	"github.com/no-src/gofs/sync"
	"github.com/no-src/gofs/version"
	"github.com/no-src/log"
)

func main() {
	parseFlags()

	// if current is subprocess, then reset the "-kill_ppid" and "-daemon"
	if config.IsSubprocess {
		config.KillPPid = false
		config.IsDaemon = false
	}

	// init logger
	var loggers []log.Logger
	loggers = append(loggers, log.NewConsoleLogger(log.Level(config.LogLevel)))
	if config.EnableFileLogger {
		filePrefix := "gofs_"
		if config.IsDaemon {
			filePrefix += "daemon_"
		}
		flogger, err := log.NewFileLoggerWithAutoFlush(log.Level(config.LogLevel), config.LogDir, filePrefix, config.LogFlush, config.LogFlushInterval)
		if err != nil {
			log.Error(err, "init file logger error")
			return
		}
		loggers = append(loggers, flogger)
	}

	log.InitDefaultLogger(log.NewMultiLogger(loggers...))
	defer log.Close()

	// print version info
	if config.PrintVersion {
		version.PrintVersion()
		return
	}

	// print about info
	if config.PrintAbout {
		about.PrintAbout()
		return
	}

	err := initFlags()
	if err != nil {
		log.Error(err, "init flags default value error")
		return
	}

	// kill parent process
	if config.KillPPid {
		daemon.KillPPid()
	}

	// start the daemon
	if config.IsDaemon {
		go signal.Notify(daemon.Shutdown)
		daemon.Daemon(config.DaemonPid, config.DaemonDelay, config.DaemonMonitorDelay)
		log.Info("daemon exited")
		return
	}

	// if enable daemon, start a worker to process the following

	userList, err := auth.ParseUsers(config.Users)
	if err != nil {
		log.Error(err, "parse users error => [%s]", config.Users)
		return
	}

	// init web server logger
	var webLogger = log.NewConsoleLogger(log.Level(config.LogLevel))
	defer webLogger.Close()
	if config.EnableFileLogger && config.EnableFileServer {
		webFileLogger, err := log.NewFileLoggerWithAutoFlush(log.Level(config.LogLevel), config.LogDir, "web_", config.LogFlush, config.LogFlushInterval)
		if err != nil {
			log.Error(err, "init the web server file logger error")
			return
		}
		webLogger = log.NewMultiLogger(webFileLogger, webLogger)
	}

	// start a file server
	if config.EnableFileServer {
		waitInit := retry.NewWaitDone()
		go func() {
			err := httpfs.StartFileServer(server.NewServerOption(config.Source, config.Dest, config.FileServerAddr, waitInit, config.EnableTLS, config.TLSCertFile, config.TLSKeyFile, userList, config.EnableFileServerCompress, webLogger, config.EnablePProf, config.PProfPrivate))
			if err != nil {
				log.Error(err, "start the file server [%s] error", config.FileServerAddr)
			}
		}()
		waitInit.Wait()
	}

	// create syncer
	syncer, err := sync.NewSync(config.Source, config.Dest, config.EnableTLS, config.TLSCertFile, config.TLSKeyFile, userList, config.EnableLogicallyDelete)
	if err != nil {
		log.Error(err, "create the instance of Sync error")
		return
	}

	// create retry
	retry := retry.NewRetry(config.RetryCount, config.RetryWait, config.RetryAsync)

	// init event log
	var eventLogger = log.NewEmptyLogger()
	defer eventLogger.Close()
	if config.EnableEventLog {
		eventFileLogger, err := log.NewFileLoggerWithAutoFlush(log.Level(config.LogLevel), config.LogDir, "event_", config.LogFlush, config.LogFlushInterval)
		if err != nil {
			log.Error(err, "init the event file logger error")
			return
		}
		eventLogger = eventFileLogger
	}

	// create monitor
	monitor, err := monitor.NewMonitor(syncer, retry, config.SyncOnce, config.EnableTLS, userList, eventLogger)
	if err != nil {
		log.Error(err, "create the instance of Monitor error")
		return
	}

	err = monitor.SyncCron(config.SyncCron)
	if err != nil {
		log.Error(err, "register sync cron task error")
		return
	}

	// start monitor
	log.Info("monitor is starting...")
	defer log.Info("gofs exited")
	go signal.Notify(monitor.Shutdown)
	defer monitor.Close()
	err = monitor.Start()
	if err != nil {
		log.Error(err, "start to monitor failed")
	}
}
