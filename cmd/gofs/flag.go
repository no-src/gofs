package main

import (
	"flag"
	"fmt"
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/conf"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/daemon"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/util"
	"github.com/no-src/log"
	"os"
	"time"
)

var (
	config conf.Config
)

func parseFlags() {

	// print help info if no arguments
	if len(os.Args) <= 1 {
		os.Args = append(os.Args, "-h")
	}

	// other
	flag.BoolVar(&config.PrintVersion, "v", false, "print the version info")
	flag.BoolVar(&config.PrintAbout, "about", false, "print the about info")

	// file sync
	core.VFSVar(&config.Source, "source", core.NewEmptyVFS(), "the source path by monitor")
	core.VFSVar(&config.Dest, "dest", core.NewEmptyVFS(), "the dest path to backup")
	flag.BoolVar(&config.SyncOnce, "sync_once", false, "sync source directory to dest directory once")
	flag.StringVar(&config.SyncCron, "sync_cron", "", "sync source directory to dest directory with cron")
	flag.BoolVar(&config.EnableLogicallyDelete, "logically_delete", false, "delete dest file logically")

	// retry
	flag.IntVar(&config.RetryCount, "retry_count", 15, "if execute failed, then retry to work -retry_count times")
	flag.DurationVar(&config.RetryWait, "retry_wait", time.Second*5, "if retry to work, wait -retry_wait time then do")
	flag.BoolVar(&config.RetryAsync, "retry_async", false, "execute retry asynchronously")

	// log
	flag.IntVar(&config.LogLevel, "log_level", int(log.InfoLevel), "set log level, default is INFO. DEBUG=0 INFO=1 WARN=2 ERROR=3")
	flag.BoolVar(&config.EnableFileLogger, "log_file", true, "enable the file logger")
	flag.StringVar(&config.LogDir, "log_dir", "./logs/", "set the directory of the log file")
	flag.BoolVar(&config.LogFlush, "log_flush", true, "enable auto flush log with interval")
	flag.DurationVar(&config.LogFlushInterval, "log_flush_interval", time.Second*3, "set the log flush interval duration, you need to enable -log_flush first")
	flag.BoolVar(&config.EnableEventLog, "log_event", false, "enable the event log")

	// daemon
	flag.BoolVar(&config.IsDaemon, "daemon", false, "enable daemon to create and monitor a subprocess to work, you can use [go build -ldflags=\"-H windowsgui\"] to build on Windows")
	flag.BoolVar(&config.DaemonPid, "daemon_pid", false, "record parent process pid, daemon process pid and worker process pid to pid file")
	flag.DurationVar(&config.DaemonDelay, "daemon_delay", time.Second, "daemon work interval, wait to create subprocess")
	flag.DurationVar(&config.DaemonMonitorDelay, "daemon_monitor_delay", time.Second*3, "daemon monitor work interval, wait to check subprocess state")
	flag.BoolVar(&config.KillPPid, "kill_ppid", false, "try to kill the parent process when it's running")
	flag.BoolVar(&config.IsSubprocess, daemon.SubprocessTag, false, "tag current process is subprocess")

	// file server
	flag.BoolVar(&config.EnableFileServer, "server", false, "start a file server to browse source directory and dest directory")
	flag.StringVar(&config.FileServerAddr, "server_addr", server.DefaultAddrHttps, "a file server binding address")
	flag.BoolVar(&config.EnableFileServerCompress, "server_compress", false, "enable response compression for the file server")
	flag.BoolVar(&config.EnablePProf, "pprof", false, "enable the pprof route")
	flag.BoolVar(&config.PProfPrivate, "pprof_private", true, "allow to access pprof route by private address and loopback address only")

	// tls transfer
	flag.BoolVar(&config.EnableTLS, "tls", true, fmt.Sprintf("enable the tls connections, if disable it, server_addr is \"%s\" default", server.DefaultAddrHttp))
	flag.StringVar(&config.TLSCertFile, "tls_cert_file", "gofs.pem", "cert file for tls connections")
	flag.StringVar(&config.TLSKeyFile, "tls_key_file", "gofs.key", "key file for tls connections")

	// login user
	flag.StringVar(&config.Users, "users", "", "the server accounts, the server allows anonymous access if there is no effective account, format like this, user1|password1,user2|password2")
	flag.IntVar(&config.RandomUserCount, "rand_user_count", 0, "the number of random server accounts, if it is greater than zero, random generate some accounts for -users")
	flag.IntVar(&config.RandomUserNameLen, "rand_user_len", 6, "the length of the random user's username")
	flag.IntVar(&config.RandomPasswordLen, "rand_pwd_len", 10, "the length of the random user's password")

	flag.Parse()
}

// initFlags init flags default value
func initFlags() error {
	if !config.EnableTLS && config.FileServerAddr == server.DefaultAddrHttps {
		config.FileServerAddr = server.DefaultAddrHttp
	}

	// if start a remote server monitor, auto enable file server
	if config.Source.Server() {
		config.EnableFileServer = true
	}

	if config.RandomUserCount > 0 && config.EnableFileServer {
		userList := auth.RandomUser(config.RandomUserCount, config.RandomUserNameLen, config.RandomPasswordLen)
		randUserStr, err := auth.ParseStringUsers(userList)
		if err != nil {
			return err
		} else {
			if len(config.Users) > 0 {
				config.Users = fmt.Sprintf("%s,%s", config.Users, randUserStr)
			} else {
				config.Users = randUserStr
			}
			log.Info("generate random users success => [%s]", config.Users)
		}
	}

	if config.EnableTLS && (config.Source.Server() || config.EnableFileServer) {
		exist, err := util.FileExist(config.TLSCertFile)
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("cert file is not found for tls => [%s], for more information, see -tls and -tls_cert_file flags", config.TLSCertFile)
		}
		exist, err = util.FileExist(config.TLSKeyFile)
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("key file is not found for tls => [%s], for more information, see -tls and -tls_key_file flags", config.TLSKeyFile)
		}
	}

	return nil
}
