package cmd

import (
	"fmt"
	"os"

	"github.com/no-src/gofs/about"
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/checksum"
	"github.com/no-src/gofs/conf"
	"github.com/no-src/gofs/daemon"
	"github.com/no-src/gofs/encrypt"
	"github.com/no-src/gofs/fs"
	"github.com/no-src/gofs/ignore"
	"github.com/no-src/gofs/internal/signal"
	"github.com/no-src/gofs/monitor"
	"github.com/no-src/gofs/result"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/server/httpfs"
	"github.com/no-src/gofs/sync"
	"github.com/no-src/gofs/util/hashutil"
	"github.com/no-src/gofs/util/httputil"
	"github.com/no-src/gofs/version"
	"github.com/no-src/gofs/wait"
	"github.com/no-src/log"
	"github.com/no-src/log/formatter"
	"github.com/no-src/log/level"
	"github.com/no-src/log/option"
)

// Run running the gofs program
func Run() result.Result {
	return RunWithArgs(os.Args)
}

// RunWithArgs running the gofs program with specified command-line arguments, starting with the program name
func RunWithArgs(args []string) result.Result {
	return RunWithConfig(parseFlags(args))
}

// RunWithConfigFile running the gofs program with specified config file
func RunWithConfigFile(path string) result.Result {
	return RunWithArgs([]string{os.Args[0], "-conf=" + path})
}

// RunWithConfig running the gofs program with specified config
func RunWithConfig(c conf.Config) result.Result {
	result := result.New()
	go runWithConfig(c, result)
	return result
}

func runWithConfig(c conf.Config, result result.Result) {
	var err error

	//  ensure all the code in this function is executed
	defer func() {
		result.DoneWithError(err)
	}()

	cp := &c
	conf.GlobalConfig = cp

	if err = parseConfigFile(cp); err != nil {
		result.InitDoneWithError(err)
		return
	}

	// if current is subprocess, then reset the "-kill_ppid" and "-daemon"
	if c.IsSubprocess {
		c.KillPPid = false
		c.IsDaemon = false
	}

	// init the default logger
	if err = initDefaultLogger(c); err != nil {
		result.InitDoneWithError(err)
		return
	}
	defer log.Close()

	var exit bool
	// execute and exit
	if exit, err = executeOnce(c); exit {
		result.InitDoneWithError(err)
		return
	}

	if exit, err = initChecksum(c); exit {
		result.InitDoneWithError(err)
		return
	}

	if err = initial(cp); err != nil {
		result.InitDoneWithError(err)
		return
	}

	// kill parent process
	if c.KillPPid {
		daemon.KillPPid()
	}

	// start the daemon
	if c.IsDaemon {
		ns := signal.Notify(daemon.Shutdown)
		go func() {
			result.RegisterNotifyHandler(ns)
		}()
		result.InitDone()
		w := wait.NewWaitDone()
		go daemon.Daemon(c.DaemonPid, c.DaemonDelay.Duration(), c.DaemonMonitorDelay.Duration(), w)
		err = w.Wait()
		return
	}

	// if enable daemon, start a worker to process the following

	userList, err := auth.ParseUsers(c.Users)
	if err != nil {
		log.Error(err, "parse users error => [%s]", c.Users)
		result.InitDoneWithError(err)
		return
	}

	// init the web server logger
	webLogger, err := initWebServerLogger(c)
	if err != nil {
		result.InitDoneWithError(err)
		return
	}
	defer webLogger.Close()

	// create retry
	r := retry.New(c.RetryCount, c.RetryWait.Duration(), c.RetryAsync)

	// start a file web server
	if err = startWebServer(c, webLogger, userList, r); err != nil {
		result.InitDoneWithError(err)
		return
	}

	// init the event log
	eventLogger, err := initEventLogger(c)
	if err != nil {
		result.InitDoneWithError(err)
		return
	}
	defer eventLogger.Close()

	// init the monitor
	m, err := initMonitor(c, userList, eventLogger, r)
	if err != nil {
		result.InitDoneWithError(err)
		return
	}

	// start monitor
	log.Info("monitor is starting...")
	defer log.Info("gofs exited")
	ns := signal.Notify(m.Shutdown)
	go func() {
		result.RegisterNotifyHandler(ns)
	}()
	defer m.Close()
	w, err := m.Start()
	result.InitDoneWithError(err)
	if err != nil {
		log.Error(err, "start to monitor failed")
		return
	}
	err = log.ErrorIf(w.Wait(), "monitor running failed")
}

func parseConfigFile(cp *conf.Config) error {
	if len(cp.Conf) > 0 {
		if err := conf.Parse(cp.Conf, cp); err != nil {
			log.Error(err, "parse config file error => [%s]", cp.Conf)
			return err
		}
	}
	return nil
}

// executeOnce execute the work and get ready to exit
func executeOnce(c conf.Config) (exit bool, err error) {
	// print version info
	if c.PrintVersion {
		version.PrintVersion("gofs")
		return true, nil
	}

	// print about info
	if c.PrintAbout {
		about.PrintAbout()
		return true, nil
	}

	// clear the deleted files
	if c.ClearDeletedPath {
		return true, log.ErrorIf(fs.ClearDeletedFile(c.Dest.Path()), "clear the deleted files error")
	}

	// decrypt the specified file or directory
	if c.Decrypt {
		return true, log.ErrorIf(encrypt.NewDecrypt(encrypt.NewOption(c)).Decrypt(), "decrypt error")
	}
	return false, nil
}

// initDefaultLogger init the default logger
func initDefaultLogger(c conf.Config) error {
	// init log formatter
	if c.LogFormat != formatter.TextFormatter {
		log.Info("switch logger format to %s", c.LogFormat)
	}
	formatter.InitDefaultFormatter(c.LogFormat)
	log.DefaultLogger().WithFormatter(formatter.New(c.LogFormat))

	var loggers []log.Logger
	loggers = append(loggers, log.NewConsoleLogger(level.Level(c.LogLevel)))
	if c.EnableFileLogger {
		filePrefix := "gofs_"
		if c.IsDaemon {
			filePrefix += "daemon_"
		}
		flogger, err := log.NewFileLoggerWithOption(option.NewFileLoggerOption(level.Level(c.LogLevel), c.LogDir, filePrefix, c.LogFlush, c.LogFlushInterval.Duration(), c.LogSplitDate))
		if err != nil {
			log.Error(err, "init file logger error")
			return err
		}
		loggers = append(loggers, flogger)
	}

	log.InitDefaultLoggerWithSample(log.NewMultiLogger(loggers...), c.LogSampleRate)
	return nil
}

// initWebServerLogger init the web server logger
func initWebServerLogger(c conf.Config) (log.Logger, error) {
	var webLogger = log.NewConsoleLogger(level.Level(c.LogLevel))
	if c.EnableFileLogger && c.EnableFileServer {
		webFileLogger, err := log.NewFileLoggerWithOption(option.NewFileLoggerOption(level.Level(c.LogLevel), c.LogDir, "web_", c.LogFlush, c.LogFlushInterval.Duration(), c.LogSplitDate))
		if err != nil {
			log.Error(err, "init the web server file logger error")
			return nil, err
		}
		webLogger = log.NewMultiLogger(webFileLogger, webLogger)
	}
	return webLogger, nil
}

// startWebServer start a file web server
func startWebServer(c conf.Config, webLogger log.Logger, userList []*auth.User, r retry.Retry) error {
	if c.EnableFileServer {
		waitInit := wait.NewWaitDone()
		go func() {
			httpfs.StartFileServer(server.NewServerOption(c, waitInit, userList, webLogger, r))
		}()

		return log.ErrorIf(waitInit.Wait(), "start the file server [%s] error", c.FileServerAddr)
	}
	return nil
}

// initEventLogger init the event logger
func initEventLogger(c conf.Config) (log.Logger, error) {
	var eventLogger = log.NewEmptyLogger()
	if c.EnableEventLog {
		eventFileLogger, err := log.NewFileLoggerWithOption(option.NewFileLoggerOption(level.Level(c.LogLevel), c.LogDir, "event_", c.LogFlush, c.LogFlushInterval.Duration(), c.LogSplitDate))
		if err != nil {
			log.Error(err, "init the event file logger error")
			return nil, err
		}
		eventLogger = eventFileLogger
	}
	return eventLogger, nil
}

// initMonitor init the monitor
func initMonitor(c conf.Config, userList []*auth.User, eventLogger log.Logger, r retry.Retry) (monitor.Monitor, error) {
	// create syncer
	syncer, err := sync.NewSync(sync.NewSyncOption(c, userList, r))
	if err != nil {
		log.Error(err, "create the instance of Sync error")
		return nil, err
	}

	// create monitor
	m, err := monitor.NewMonitor(monitor.NewMonitorOption(c, syncer, r, userList, eventLogger))
	if err != nil {
		log.Error(err, "create the instance of Monitor error")
		return nil, err
	}

	err = m.SyncCron(c.SyncCron)
	if err != nil {
		log.Error(err, "register sync cron task error")
		return nil, err
	}
	return m, nil
}

func initial(cp *conf.Config) (err error) {
	// init ignore config
	if err = log.ErrorIf(ignore.Init(cp.IgnoreConf, cp.IgnoreDeletedPath), "init ignore config error"); err != nil {
		return err
	}

	if err = log.ErrorIf(initDefaultValue(cp), "init default value of config error"); err != nil {
		return err
	}

	// init default http util
	return log.ErrorIf(httputil.Init(cp.TLSInsecureSkipVerify, cp.TLSCertFile), "init http util error")
}

func initChecksum(c conf.Config) (exit bool, err error) {
	// init default hash algorithm
	if err = log.ErrorIf(hashutil.InitDefaultHash(c.ChecksumAlgorithm), "init default hash algorithm error"); err != nil {
		return true, err
	}

	// calculate checksum
	if c.Checksum {
		return true, checksum.PrintChecksum(c.Source.Path(), c.ChunkSize, c.CheckpointCount)
	}
	return false, nil
}

// initDefaultValue init default value of config
func initDefaultValue(cp *conf.Config) error {
	initFileServer(cp)

	if err := generateRandomUser(cp); err != nil {
		return err
	}

	if err := checkTLS(*cp); err != nil {
		return err
	}

	return nil
}

// initFileServer init config about the file server
func initFileServer(cp *conf.Config) {
	if !cp.EnableTLS && cp.FileServerAddr == server.DefaultAddrHttps {
		cp.FileServerAddr = server.DefaultAddrHttp
	}

	// if start a remote server monitor, auto enable file server
	if cp.Source.Server() {
		cp.EnableFileServer = true
	}
}

// generateRandomUser check and generate some random user
func generateRandomUser(cp *conf.Config) error {
	if cp.RandomUserCount > 0 && cp.EnableFileServer {
		userList, err := auth.RandomUser(cp.RandomUserCount, cp.RandomUserNameLen, cp.RandomPasswordLen, cp.RandomDefaultPerm)
		if err != nil {
			return err
		}
		randUserStr, err := auth.ParseStringUsers(userList)
		if err != nil {
			return err
		}
		if len(cp.Users) > 0 {
			cp.Users = fmt.Sprintf("%s,%s", cp.Users, randUserStr)
		} else {
			cp.Users = randUserStr
		}
		log.Info("generate random users success => [%s]", cp.Users)
	}
	return nil
}

// checkTLS check cert and key file of the TLS
func checkTLS(c conf.Config) error {
	if c.EnableTLS && (c.Source.Server() || c.EnableFileServer) {
		exist, err := fs.FileExist(c.TLSCertFile)
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("cert file is not found for tls => [%s], for more information, see -tls and -tls_cert_file flags", c.TLSCertFile)
		}
		exist, err = fs.FileExist(c.TLSKeyFile)
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("key file is not found for tls => [%s], for more information, see -tls and -tls_key_file flags", c.TLSKeyFile)
		}
	}
	return nil
}
