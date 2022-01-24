package main

import (
	"flag"
	"fmt"
	"github.com/no-src/gofs/auth"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/daemon"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/util"
	"github.com/no-src/log"
	"os"
	"time"
)

var (

	// other
	printVersion bool

	// file sync
	source                core.VFS
	dest                  core.VFS
	syncOnce              bool
	syncCron              string
	enableLogicallyDelete bool

	// retry
	retryCount int
	retryWait  time.Duration
	retryAsync bool

	// log
	logLevel         int
	enableFileLogger bool
	logDir           string
	logFlush         bool
	logFlushInterval time.Duration
	enableEventLog   bool

	// daemon
	isDaemon           bool
	daemonPid          bool
	daemonDelay        time.Duration
	daemonMonitorDelay time.Duration
	killPPid           bool
	isSubprocess       bool

	// file server
	enableFileServer         bool
	fileServerAddr           string
	enableFileServerCompress bool
	enablePprof              bool
	pprofPrivate             bool

	// tls transfer
	enableTLS   bool
	tlsCertFile string
	tlsKeyFile  string

	// login user
	users             string
	randomUserCount   int
	randomUserNameLen int
	randomPasswordLen int
)

func parseFlags() {

	// print help info if no arguments
	if len(os.Args) <= 1 {
		os.Args = append(os.Args, "-h")
	}

	// other
	flag.BoolVar(&printVersion, "v", false, "print the version info")

	// file sync
	core.VFSVar(&source, "source", core.NewEmptyVFS(), "the source path by monitor")
	core.VFSVar(&dest, "dest", core.NewEmptyVFS(), "the dest path to backup")
	flag.BoolVar(&syncOnce, "sync_once", false, "sync source directory to dest directory once")
	flag.StringVar(&syncCron, "sync_cron", "", "sync source directory to dest directory with cron")
	flag.BoolVar(&enableLogicallyDelete, "logically_delete", false, "delete dest file logically")

	// retry
	flag.IntVar(&retryCount, "retry_count", 15, "if execute failed, then retry to work -retry_count times")
	flag.DurationVar(&retryWait, "retry_wait", time.Second*5, "if retry to work, wait -retry_wait time then do")
	flag.BoolVar(&retryAsync, "retry_async", false, "execute retry asynchronously")

	// log
	flag.IntVar(&logLevel, "log_level", int(log.InfoLevel), "set log level, default is INFO. DEBUG=0 INFO=1 WARN=2 ERROR=3")
	flag.BoolVar(&enableFileLogger, "log_file", true, "enable the file logger")
	flag.StringVar(&logDir, "log_dir", "./logs/", "set the directory of the log file")
	flag.BoolVar(&logFlush, "log_flush", true, "enable auto flush log with interval")
	flag.DurationVar(&logFlushInterval, "log_flush_interval", time.Second*3, "set the log flush interval duration, you need to enable -log_flush first")
	flag.BoolVar(&enableEventLog, "log_event", false, "enable the event log")

	// daemon
	flag.BoolVar(&isDaemon, "daemon", false, "enable daemon to create and monitor a subprocess to work, you can use [go build -ldflags=\"-H windowsgui\"] to build on Windows")
	flag.BoolVar(&daemonPid, "daemon_pid", false, "record parent process pid, daemon process pid and worker process pid to pid file")
	flag.DurationVar(&daemonDelay, "daemon_delay", time.Second, "daemon work interval, wait to create subprocess")
	flag.DurationVar(&daemonMonitorDelay, "daemon_monitor_delay", time.Second*3, "daemon monitor work interval, wait to check subprocess state")
	flag.BoolVar(&killPPid, "kill_ppid", false, "try to kill the parent process when it's running")
	flag.BoolVar(&isSubprocess, daemon.SubprocessTag, false, "tag current process is subprocess")

	// file server
	flag.BoolVar(&enableFileServer, "server", false, "start a file server to browse source directory and dest directory")
	flag.StringVar(&fileServerAddr, "server_addr", server.DefaultAddrHttps, "a file server binding address")
	flag.BoolVar(&enableFileServerCompress, "server_compress", false, "enable response compression for the file server")
	flag.BoolVar(&enablePprof, "pprof", false, "enable the pprof route")
	flag.BoolVar(&pprofPrivate, "pprof_private", true, "allow to access pprof route by private address and loopback address only")

	// tls transfer
	flag.BoolVar(&enableTLS, "tls", true, fmt.Sprintf("enable the tls connections, if disable it, server_addr is \"%s\" default", server.DefaultAddrHttp))
	flag.StringVar(&tlsCertFile, "tls_cert_file", "gofs.pem", "cert file for tls connections")
	flag.StringVar(&tlsKeyFile, "tls_key_file", "gofs.key", "key file for tls connections")

	// login user
	flag.StringVar(&users, "users", "", "the server accounts, the server allows anonymous access if there is no effective account, format like this, user1|password1,user2|password2")
	flag.IntVar(&randomUserCount, "rand_user_count", 0, "the number of random server accounts, if it is greater than zero, random generate some accounts for -users")
	flag.IntVar(&randomUserNameLen, "rand_user_len", 6, "the length of the random user's username")
	flag.IntVar(&randomPasswordLen, "rand_pwd_len", 10, "the length of the random user's password")

	flag.Parse()
}

// initFlags init flags default value
func initFlags() error {
	if !enableTLS && fileServerAddr == server.DefaultAddrHttps {
		fileServerAddr = server.DefaultAddrHttp
	}

	// if start a remote server monitor, auto enable file server
	if source.Server() {
		enableFileServer = true
	}

	if randomUserCount > 0 && enableFileServer {
		userList := auth.RandomUser(randomUserCount, randomUserNameLen, randomPasswordLen)
		randUserStr, err := auth.ParseStringUsers(userList)
		if err != nil {
			return err
		} else {
			if len(users) > 0 {
				users = fmt.Sprintf("%s,%s", users, randUserStr)
			} else {
				users = randUserStr
			}
			log.Info("generate random users success => [%s]", users)
		}
	}

	if enableTLS && (source.Server() || enableFileServer) {
		exist, err := util.FileExist(tlsCertFile)
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("cert file is not found for tls => [%s], for more information, see -tls and -tls_cert_file flags", tlsCertFile)
		}
		exist, err = util.FileExist(tlsKeyFile)
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("key file is not found for tls => [%s], for more information, see -tls and -tls_key_file flags", tlsKeyFile)
		}
	}

	return nil
}
