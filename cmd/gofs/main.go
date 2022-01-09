package main

import (
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/daemon"
	"github.com/no-src/gofs/monitor"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/server/fs"
	"github.com/no-src/gofs/sync"
	"github.com/no-src/gofs/version"
	"github.com/no-src/log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	parseFlags()

	// if current is subprocess, then reset the "-kill_ppid" and "-daemon"
	if isSubprocess {
		killPPid = false
		isDaemon = false
	}

	// init logger
	var loggers []log.Logger
	loggers = append(loggers, log.NewConsoleLogger(log.Level(logLevel)))
	if enableFileLogger {
		filePrefix := "gofs_"
		if isDaemon {
			filePrefix += "daemon_"
		}
		flogger, err := log.NewFileLoggerWithAutoFlush(log.Level(logLevel), logDir, filePrefix, logFlush, logFlushInterval)
		if err != nil {
			log.Error(err, "init file logger error")
			return
		}
		loggers = append(loggers, flogger)
	}

	log.InitDefaultLogger(log.NewMultiLogger(loggers...))
	defer log.Close()

	// print version info
	if printVersion {
		version.PrintVersionInfo()
		return
	}

	err := initFlags()
	if err != nil {
		log.Error(err, "init flags default value error")
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

	userList, err := auth.ParseUsers(users)
	if err != nil {
		log.Error(err, "parse users error => [%s]", users)
		return
	}

	// init web server logger
	var webLogger = log.NewConsoleLogger(log.Level(logLevel))
	defer webLogger.Close()
	if enableFileLogger && enableFileServer {
		webFileLogger, err := log.NewFileLoggerWithAutoFlush(log.Level(logLevel), logDir, "web_", logFlush, logFlushInterval)
		if err != nil {
			log.Error(err, "init the web server file logger error")
			return
		}
		webLogger = log.NewMultiLogger(webFileLogger, webLogger)
	}

	// start a file server
	if enableFileServer {
		waitInit := retry.NewWaitDone()
		go func() {
			err := fs.StartFileServer(server.NewServerOption(sourceVFS, targetVFS, fileServerAddr, waitInit, enableTLS, tlsCertFile, tlsKeyFile, userList, enableFileServerCompress, webLogger, enablePprof, pprofPrivate))
			if err != nil {
				log.Error(err, "start the file server [%s] error", fileServerAddr)
			}
		}()
		waitInit.Wait()
	}

	// create syncer
	syncer, err := sync.NewSync(sourceVFS, targetVFS, bufSize, enableTLS, tlsCertFile, tlsKeyFile, userList)
	if err != nil {
		log.Error(err, "create the instance of Sync error")
		return
	}

	// create retry
	retry := retry.NewRetry(retryCount, retryWait, retryAsync)

	// create monitor
	monitor, err := monitor.NewMonitor(syncer, retry, syncOnce, enableTLS, userList)
	if err != nil {
		log.Error(err, "create the instance of Monitor error")
		return
	}

	err = monitor.SyncCron(syncCron)
	if err != nil {
		log.Error(err, "register sync cron task error")
		return
	}

	// start monitor
	log.Info("monitor is starting...")
	defer log.Info("gofs exited!")
	go notify(monitor.Shutdown)
	defer monitor.Close()
	err = monitor.Start()
	if err != nil {
		log.Error(err, "start to monitor failed")
	}
}

func notify(shutdown func() error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGTERM)
	for {
		s := <-c
		switch s {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGKILL, syscall.SIGTERM:
			log.Debug("received a signal [%s], waiting to exit", s.String())
			err := shutdown()
			if err != nil {
				log.Error(err, "shutdown error")
				break
			} else {
				signal.Stop(c)
				close(c)
				return
			}
		default:
			log.Debug("received a signal [%s], ignore it", s.String())
			break
		}
	}
}
