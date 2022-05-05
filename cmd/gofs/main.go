package main

import (
	"github.com/no-src/gofs/about"
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/checksum"
	"github.com/no-src/gofs/conf"
	"github.com/no-src/gofs/daemon"
	"github.com/no-src/gofs/fs"
	"github.com/no-src/gofs/ignore"
	"github.com/no-src/gofs/internal/signal"
	"github.com/no-src/gofs/monitor"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/server/httpfs"
	"github.com/no-src/gofs/sync"
	"github.com/no-src/gofs/version"
	"github.com/no-src/gofs/wait"
	"github.com/no-src/log"
)

func main() {
	parseFlags()

	if parseConfig() != nil {
		return
	}

	// if current is subprocess, then reset the "-kill_ppid" and "-daemon"
	if config.IsSubprocess {
		config.KillPPid = false
		config.IsDaemon = false
	}

	// init the default logger
	if initDefaultLogger() != nil {
		return
	}
	defer log.Close()

	// execute and exit
	if executeOnce() {
		return
	}

	// init ignore config
	if err := ignore.Init(config.IgnoreConf, config.IgnoreDeletedPath); err != nil {
		log.Error(err, "init ignore config error")
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
		daemon.Daemon(config.DaemonPid, config.DaemonDelay.Duration(), config.DaemonMonitorDelay.Duration())
		log.Info("daemon exited")
		return
	}

	// if enable daemon, start a worker to process the following

	userList, err := auth.ParseUsers(config.Users)
	if err != nil {
		log.Error(err, "parse users error => [%s]", config.Users)
		return
	}

	// init the web server logger
	webLogger, err := initWebServerLogger()
	if err != nil {
		return
	}
	defer webLogger.Close()

	// start a file web server
	startWebServer(webLogger, userList)

	// init the event log
	eventLogger, err := initEventLogger()
	if err != nil {
		return
	}
	defer eventLogger.Close()

	// init the monitor
	m, err := initMonitor(userList, eventLogger)
	if err != nil {
		return
	}

	// start monitor
	log.Info("monitor is starting...")
	defer log.Info("gofs exited")
	go signal.Notify(m.Shutdown)
	defer m.Close()
	err = m.Start()
	if err != nil {
		log.Error(err, "start to monitor failed")
	}
}

func parseConfig() error {
	if len(config.Conf) > 0 {
		if err := conf.Parse(config.Conf, &config); err != nil {
			log.Error(err, "parse config file error => [%s]", config.Conf)
			return err
		}
	}
	return nil
}

// executeOnce execute the work and get ready to exit
func executeOnce() (exit bool) {
	// print version info
	if config.PrintVersion {
		version.PrintVersion()
		return true
	}

	// print about info
	if config.PrintAbout {
		about.PrintAbout()
		return true
	}

	// clear the deleted files
	if config.ClearDeletedPath {
		log.ErrorIf(fs.ClearDeletedFile(config.Dest.Path()), "clear the deleted files error")
		return true
	}

	// calculate checksum
	if config.Checksum {
		checksum.PrintChecksum(config.Source.Path(), config.ChunkSize, config.CheckpointCount)
		return true
	}

	return false
}

// initDefaultLogger init the default logger
func initDefaultLogger() error {
	var loggers []log.Logger
	loggers = append(loggers, log.NewConsoleLogger(log.Level(config.LogLevel)))
	if config.EnableFileLogger {
		filePrefix := "gofs_"
		if config.IsDaemon {
			filePrefix += "daemon_"
		}
		flogger, err := log.NewFileLoggerWithAutoFlush(log.Level(config.LogLevel), config.LogDir, filePrefix, config.LogFlush, config.LogFlushInterval.Duration())
		if err != nil {
			log.Error(err, "init file logger error")
			return err
		}
		loggers = append(loggers, flogger)
	}

	log.InitDefaultLogger(log.NewMultiLogger(loggers...))
	return nil
}

// initWebServerLogger init the web server logger
func initWebServerLogger() (log.Logger, error) {
	var webLogger = log.NewConsoleLogger(log.Level(config.LogLevel))
	if config.EnableFileLogger && config.EnableFileServer {
		webFileLogger, err := log.NewFileLoggerWithAutoFlush(log.Level(config.LogLevel), config.LogDir, "web_", config.LogFlush, config.LogFlushInterval.Duration())
		if err != nil {
			log.Error(err, "init the web server file logger error")
			return nil, err
		}
		webLogger = log.NewMultiLogger(webFileLogger, webLogger)
	}
	return webLogger, nil
}

// startWebServer start a file web server
func startWebServer(webLogger log.Logger, userList []*auth.User) {
	if config.EnableFileServer {
		waitInit := wait.NewWaitDone()
		go func() {
			log.ErrorIf(httpfs.StartFileServer(server.NewServerOption(config, waitInit, userList, webLogger)), "start the file server [%s] error", config.FileServerAddr)
		}()
		waitInit.Wait()
	}
}

// initEventLogger init the event logger
func initEventLogger() (log.Logger, error) {
	var eventLogger = log.NewEmptyLogger()
	if config.EnableEventLog {
		eventFileLogger, err := log.NewFileLoggerWithAutoFlush(log.Level(config.LogLevel), config.LogDir, "event_", config.LogFlush, config.LogFlushInterval.Duration())
		if err != nil {
			log.Error(err, "init the event file logger error")
			return nil, err
		}
		eventLogger = eventFileLogger
	}
	return eventLogger, nil
}

// initMonitor init the monitor
func initMonitor(userList []*auth.User, eventLogger log.Logger) (monitor.Monitor, error) {
	// create syncer
	syncer, err := sync.NewSync(sync.NewSyncOption(config, userList))
	if err != nil {
		log.Error(err, "create the instance of Sync error")
		return nil, err
	}

	// create retry
	r := retry.New(config.RetryCount, config.RetryWait.Duration(), config.RetryAsync)

	// create monitor
	m, err := monitor.NewMonitor(syncer, r, config.SyncOnce, config.EnableTLS, userList, eventLogger)
	if err != nil {
		log.Error(err, "create the instance of Monitor error")
		return nil, err
	}

	err = m.SyncCron(config.SyncCron)
	if err != nil {
		log.Error(err, "register sync cron task error")
		return nil, err
	}
	return m, nil
}
