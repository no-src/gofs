package main

import (
	"flag"
	"github.com/no-src/gofs/daemon"
	"github.com/no-src/log"
	"time"
)

var (
	SrcPath            string
	TargetPath         string
	LogLevel           int
	FileLogger         bool
	LogDir             string
	LogFlush           bool
	LogFlushInterval   time.Duration
	RetryCount         int
	RetryWait          time.Duration
	BufSize            int
	PrintVersion       bool
	SyncOnce           bool
	Daemon             bool
	DaemonPid          bool
	DaemonDelay        time.Duration
	DaemonMonitorDelay time.Duration
	KillPPid           bool
	IsSubprocess       bool
	FileServer         bool
	FileServerAddr     string
)

func parseFlags() {
	flag.BoolVar(&PrintVersion, "v", false, "print version info")
	flag.StringVar(&SrcPath, "src", "", "source path by monitor")
	flag.StringVar(&TargetPath, "target", "", "target path to backup")
	flag.IntVar(&LogLevel, "log_level", int(log.InfoLevel), "set log level, default is INFO. DEBUG=0 INFO=1 WARN=2 ERROR=3")
	flag.BoolVar(&FileLogger, "log_file", false, "enable file logger")
	flag.StringVar(&LogDir, "log_dir", "./logs/", "set log file's dir")
	flag.BoolVar(&LogFlush, "log_flush", false, "enable auto flush log with interval")
	flag.DurationVar(&LogFlushInterval, "log_flush_interval", time.Second*3, "set log flush interval duration, you need to enable log_flush first")
	flag.IntVar(&RetryCount, "retry_count", 15, "if execute failed, then retry to work retry_count times")
	flag.DurationVar(&RetryWait, "retry_wait", time.Second*5, "if retry to work, wait retry_wait time then do")
	flag.IntVar(&BufSize, "buf_size", 1024*1024, "read and write buffer byte size")
	flag.BoolVar(&SyncOnce, "sync_once", false, "sync src directory to target directory once")
	flag.BoolVar(&Daemon, "daemon", false, "enable daemon to create and monitor a subprocess to work, you can use [go build -ldflags=\"-H windowsgui\"] to build on Windows")
	flag.BoolVar(&DaemonPid, "daemon_pid", false, "record parent process pid, daemon process pid and worker process pid to pid file")
	flag.DurationVar(&DaemonDelay, "daemon_delay", time.Second, "daemon work interval, wait to create subprocess")
	flag.DurationVar(&DaemonMonitorDelay, "daemon_monitor_delay", time.Second*3, "daemon monitor work interval, wait to check subprocess state")
	flag.BoolVar(&KillPPid, "kill_ppid", false, "try to kill the parent process when it's running")
	flag.BoolVar(&IsSubprocess, daemon.SubprocessTag, false, "tag current process is subprocess")
	flag.BoolVar(&FileServer, "server", false, "start a file server to browse source directory and target directory")
	flag.StringVar(&FileServerAddr, "server_addr", ":9015", "a file server binding address")
	flag.Parse()
}
