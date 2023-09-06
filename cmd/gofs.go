package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/checksum"
	"github.com/no-src/gofs/conf"
	"github.com/no-src/gofs/daemon"
	"github.com/no-src/gofs/encrypt"
	"github.com/no-src/gofs/flag"
	"github.com/no-src/gofs/fs"
	"github.com/no-src/gofs/ignore"
	"github.com/no-src/gofs/internal/about"
	"github.com/no-src/gofs/internal/logger"
	"github.com/no-src/gofs/internal/signal"
	"github.com/no-src/gofs/internal/version"
	"github.com/no-src/gofs/monitor"
	"github.com/no-src/gofs/report"
	"github.com/no-src/gofs/result"
	"github.com/no-src/gofs/retry"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/server/httpfs"
	"github.com/no-src/gofs/sync"
	"github.com/no-src/gofs/wait"
	"github.com/no-src/log/level"
)

// Run running the gofs program
func Run() result.Result {
	return RunWithArgs(os.Args)
}

// RunWithArgs running the gofs program with specified command-line arguments, starting with the program name
func RunWithArgs(args []string) result.Result {
	return RunWithConfig(flag.ParseFlags(args))
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

// RunWithConfigContent running the gofs program with specified config content
func RunWithConfigContent(content string, ext string) result.Result {
	var c conf.Config
	err := conf.ParseContent([]byte(content), ext, &c)
	if err != nil {
		result := result.New()
		result.InitDoneWithError(err)
		result.DoneWithError(err)
		return result
	}
	return RunWithConfig(c)
}

//gocyclo:ignore
func runWithConfig(c conf.Config, result result.Result) {
	var err error

	//  ensure all the code in this function is executed
	defer func() {
		result.DoneWithError(err)
	}()

	cp := &c

	if err = parseConfigFile(cp); err != nil {
		result.InitDoneWithError(err)
		return
	}

	// if current is subprocess, then reset the "-kill_ppid" and "-daemon"
	if c.IsSubprocess {
		c.KillPPid = false
		c.IsDaemon = false
	}

	switchDebug := false
	if c.DryRun && c.LogLevel != int(level.DebugLevel) {
		c.LogLevel = int(level.DebugLevel)
		switchDebug = true
	}

	// init the default logger
	var logger *logger.Logger
	if logger, err = initDefaultLogger(c); err != nil {
		result.InitDoneWithError(err)
		return
	}
	defer logger.Close()

	if switchDebug {
		logger.Info("to be able to see more details, force enable the log with debug level in dry run mode!")
	}

	var exit bool
	// execute and exit
	if exit, err = executeOnce(c, logger); exit {
		result.InitDoneWithError(err)
		return
	}

	if err = initDefaultValue(cp, logger); err != nil {
		logger.Error(err, "init default value of config error")
		result.InitDoneWithError(err)
		return
	}

	// kill parent process
	daemon := daemon.New(logger)
	if c.KillPPid {
		daemon.KillPPid()
	}

	// start the daemon
	if c.IsDaemon {
		var args []string
		args, err = c.ToArgs()
		if err != nil {
			result.InitDoneWithError(err)
			return
		}

		ns, ss := signal.Notify(daemon.Shutdown, logger)
		go func() {
			result.RegisterNotifyHandler(ns)
		}()
		result.InitDone()
		w := wait.NewWaitDone()
		go daemon.Run(args, c.DaemonPid, c.DaemonDelay.Duration(), c.DaemonMonitorDelay.Duration(), w)
		err = w.Wait()
		ss()
		return
	}

	// if enable daemon, start a worker to process the following

	userList, err := auth.ParseUsers(c.Users)
	if err != nil {
		logger.Error(err, "parse users error => [%s]", c.Users)
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

	reporter := report.NewReporter()
	// start a file web server
	if err = startWebServer(c, webLogger, userList, r, reporter, logger); err != nil {
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

	pi, err := ignore.NewPathIgnore(c.IgnoreConf, c.IgnoreDeletedPath)
	if err != nil {
		logger.Error(err, "init ignore config error")
		result.InitDoneWithError(err)
		return
	}

	// init the monitor
	m, err := initMonitor(c, userList, eventLogger, r, pi, reporter, logger)
	if err != nil {
		result.InitDoneWithError(err)
		return
	}

	// start monitor
	logger.Info("monitor is starting...")
	defer logger.Info("gofs exited")
	ns, ss := signal.Notify(m.Shutdown, logger)
	go func() {
		result.RegisterNotifyHandler(ns)
	}()
	defer m.Close()
	w, err := m.Start()
	result.InitDoneWithError(err)
	if err != nil {
		logger.Error(err, "start to monitor failed")
		return
	}
	err = logger.ErrorIf(w.Wait(), "monitor running failed")
	ss()
}

func parseConfigFile(cp *conf.Config) error {
	if len(cp.Conf) > 0 {
		if err := conf.Parse(cp.Conf, cp); err != nil {
			innerLogger.Error(err, "parse config file error => [%s]", cp.Conf)
			return err
		}
	}
	return nil
}

// executeOnce execute the work and get ready to exit
func executeOnce(c conf.Config, logger *logger.Logger) (exit bool, err error) {
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
		return true, logger.ErrorIf(fs.ClearDeletedFile(c.Dest.Path()), "clear the deleted files error")
	}

	// decrypt the specified file or directory
	if c.Decrypt {
		dec, err := encrypt.NewDecrypt(encrypt.NewOption(c))
		if err != nil {
			logger.Error(err, "init decrypt component error")
			return true, err
		}
		return true, logger.ErrorIf(dec.Decrypt(), "decrypt error")
	}

	// calculate checksum
	if c.Checksum {
		return true, checksum.PrintChecksum(c.Source.Path(), c.ChunkSize, c.CheckpointCount, c.ChecksumAlgorithm, logger)
	}
	return false, nil
}

// startWebServer start a file web server
func startWebServer(c conf.Config, webLogger *logger.Logger, userList []*auth.User, r retry.Retry, reporter report.Reporter, logger *logger.Logger) error {
	if c.EnableFileServer {
		waitInit := wait.NewWaitDone()
		go func() {
			httpfs.StartFileServer(server.NewServerOption(c, waitInit, userList, webLogger, r, reporter))
		}()
		return logger.ErrorIf(waitInit.Wait(), "start the file server [%s] error", c.FileServerAddr)
	}
	return nil
}

// initMonitor init the monitor
func initMonitor(c conf.Config, userList []*auth.User, eventWriter io.Writer, r retry.Retry, pi ignore.PathIgnore, reporter report.Reporter, logger *logger.Logger) (monitor.Monitor, error) {
	// create syncer
	syncer, err := sync.NewSync(sync.NewSyncOption(c, userList, r, pi, reporter, logger))
	if err != nil {
		logger.Error(err, "create the instance of Sync error")
		return nil, err
	}

	// create monitor
	m, err := monitor.NewMonitor(monitor.NewMonitorOption(c, syncer, r, userList, eventWriter, pi, reporter, logger), RunWithConfigContent)
	if err != nil {
		logger.Error(err, "create the instance of Monitor error")
		return nil, err
	}

	err = m.SyncCron(c.SyncCron)
	if err != nil {
		logger.Error(err, "register sync cron task error")
		return nil, err
	}
	return m, nil
}

// initDefaultValue init default value of config
func initDefaultValue(cp *conf.Config, logger *logger.Logger) error {
	initFileServer(cp)

	if err := generateRandomUser(cp, logger); err != nil {
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
func generateRandomUser(cp *conf.Config, logger *logger.Logger) error {
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
		logger.Info("generate random users success => [%s]", cp.Users)
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
