package main

import (
	"flag"
	"fmt"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/daemon"
	"github.com/no-src/gofs/server"
	"github.com/no-src/log"
	"os"
	"time"
)

var (
	sourceVFS          core.VFS
	targetVFS          core.VFS
	logLevel           int
	fileLogger         bool
	logDir             string
	logFlush           bool
	logFlushInterval   time.Duration
	retryCount         int
	retryWait          time.Duration
	retryAsync         bool
	bufSize            int
	printVersion       bool
	syncOnce           bool
	isDaemon           bool
	daemonPid          bool
	daemonDelay        time.Duration
	daemonMonitorDelay time.Duration
	killPPid           bool
	isSubprocess       bool
	fileServer         bool
	fileServerAddr     string
	fileServerTLS      bool
	certFile           string
	keyFile            string
	fileServerUsers    string
	fileServerTemplate string
)

func parseFlags() {

	// print help info if no arguments
	if len(os.Args) <= 1 {
		os.Args = append(os.Args, "-h")
	}

	flag.BoolVar(&printVersion, "v", false, "print version info")
	core.VFSVar(&sourceVFS, "src", core.NewEmptyVFS(), "source path by monitor")
	core.VFSVar(&targetVFS, "target", core.NewEmptyVFS(), "target path to backup")
	flag.IntVar(&logLevel, "log_level", int(log.InfoLevel), "set log level, default is INFO. DEBUG=0 INFO=1 WARN=2 ERROR=3")
	flag.BoolVar(&fileLogger, "log_file", false, "enable file logger")
	flag.StringVar(&logDir, "log_dir", "./logs/", "set log file's dir")
	flag.BoolVar(&logFlush, "log_flush", false, "enable auto flush log with interval")
	flag.DurationVar(&logFlushInterval, "log_flush_interval", time.Second*3, "set log flush interval duration, you need to enable log_flush first")
	flag.IntVar(&retryCount, "retry_count", 15, "if execute failed, then retry to work retry_count times")
	flag.DurationVar(&retryWait, "retry_wait", time.Second*5, "if retry to work, wait retry_wait time then do")
	flag.BoolVar(&retryAsync, "retry_async", false, "execute retry asynchronously")
	flag.IntVar(&bufSize, "buf_size", 1024*1024, "read and write buffer byte size")
	flag.BoolVar(&syncOnce, "sync_once", false, "sync src directory to target directory once")
	flag.BoolVar(&isDaemon, "daemon", false, "enable daemon to create and monitor a subprocess to work, you can use [go build -ldflags=\"-H windowsgui\"] to build on Windows")
	flag.BoolVar(&daemonPid, "daemon_pid", false, "record parent process pid, daemon process pid and worker process pid to pid file")
	flag.DurationVar(&daemonDelay, "daemon_delay", time.Second, "daemon work interval, wait to create subprocess")
	flag.DurationVar(&daemonMonitorDelay, "daemon_monitor_delay", time.Second*3, "daemon monitor work interval, wait to check subprocess state")
	flag.BoolVar(&killPPid, "kill_ppid", false, "try to kill the parent process when it's running")
	flag.BoolVar(&isSubprocess, daemon.SubprocessTag, false, "tag current process is subprocess")
	flag.BoolVar(&fileServer, "server", false, "start a file server to browse source directory and target directory")
	flag.StringVar(&fileServerAddr, "server_addr", server.DefaultAddrHttps, "a file server binding address")
	flag.BoolVar(&fileServerTLS, "server_tls", true, fmt.Sprintf("enable https for file server, if disable it, server_addr is \"%s\" default", server.DefaultAddrHttp))
	flag.StringVar(&fileServerUsers, "server_users", "", "the file server accounts, the file server allows anonymous access if there is no effective account, format like this, user1|password1,user2|password2")
	flag.StringVar(&fileServerTemplate, "server_tmpl", "./template/*.html", "the file server template pattern")
	flag.StringVar(&certFile, "tls_cert_file", "gofs.pem", "cert file for https connections")
	flag.StringVar(&keyFile, "tls_key_file", "gofs.key", "key file for https connections")
	flag.Parse()

	if !fileServerTLS && fileServerAddr == server.DefaultAddrHttps {
		fileServerAddr = server.DefaultAddrHttp
	}
}
