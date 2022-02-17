package conf

import (
	"github.com/no-src/gofs/core"
	"time"
)

type Config struct {
	// other
	PrintVersion bool
	PrintAbout   bool

	// file sync
	Source                core.VFS
	Dest                  core.VFS
	SyncOnce              bool
	SyncCron              string
	EnableLogicallyDelete bool
	ClearDeletedPath      bool
	IgnoreConf            string
	IgnoreDeletedPath     bool

	// retry
	RetryCount int
	RetryWait  time.Duration
	RetryAsync bool

	// log
	LogLevel         int
	EnableFileLogger bool
	LogDir           string
	LogFlush         bool
	LogFlushInterval time.Duration
	EnableEventLog   bool

	// daemon
	IsDaemon           bool
	DaemonPid          bool
	DaemonDelay        time.Duration
	DaemonMonitorDelay time.Duration
	KillPPid           bool
	IsSubprocess       bool

	// file server
	EnableFileServer         bool
	FileServerAddr           string
	EnableFileServerCompress bool
	EnablePProf              bool
	PProfPrivate             bool
	EnablePushServer         bool

	// tls transfer
	EnableTLS   bool
	TLSCertFile string
	TLSKeyFile  string

	// login user
	Users             string
	RandomUserCount   int
	RandomUserNameLen int
	RandomPasswordLen int
	RandomDefaultPerm string
}
