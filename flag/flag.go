package flag

import (
	"flag"
	"fmt"
	"time"

	"github.com/no-src/gofs/conf"
	"github.com/no-src/gofs/core"
	"github.com/no-src/gofs/daemon"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/util/hashutil"
	"github.com/no-src/log/formatter"
	"github.com/no-src/log/level"
)

// ParseFlags parse config from arguments
func ParseFlags(args []string) (config conf.Config) {
	if len(args) < 1 {
		panic("at least one argument is required, starting with the program name")
	}

	// print help info if no arguments
	if len(args) <= 1 {
		args = append(args, "-h")
	}

	var cl = core.NewFlagSet(args[0], flag.ExitOnError)

	// other
	cl.BoolVar(&config.PrintVersion, "v", false, "print the version info")
	cl.BoolVar(&config.PrintAbout, "about", false, "print the about info")
	cl.StringVar(&config.Conf, "conf", "", "the path of config file")

	// file sync
	cl.VFSVar(&config.Source, "source", core.NewEmptyVFS(), "the source path by monitor")
	cl.VFSVar(&config.Dest, "dest", core.NewEmptyVFS(), "the dest path to backup")
	cl.BoolVar(&config.SyncOnce, "sync_once", false, "sync source directory to dest directory once")
	cl.StringVar(&config.SyncCron, "sync_cron", "", "sync source directory to dest directory with cron")
	cl.BoolVar(&config.EnableLogicallyDelete, "logically_delete", false, "delete dest file logically")
	cl.BoolVar(&config.ClearDeletedPath, "clear_deleted", false, "remove all of the deleted files in the dest path")
	cl.StringVar(&config.IgnoreConf, "ignore_conf", "", "a config file of the ignore component")
	cl.BoolVar(&config.IgnoreDeletedPath, "ignore_deleted", true, "ignore to sync the deleted file")
	cl.Int64Var(&config.ChunkSize, "chunk_size", 1024*1024, "the chunk size of the big file")
	cl.IntVar(&config.CheckpointCount, "checkpoint_count", 10, "use the checkpoint in the file to reduce transfer unmodified file chunks")
	cl.BoolVar(&config.ForceChecksum, "force_checksum", false, "if the file size and file modification time of the source file is equal to the destination file and -force_checksum is false, then ignore the current file transfer")
	cl.StringVar(&config.ChecksumAlgorithm, "checksum_algorithm", hashutil.DefaultHash, "set the default hash algorithm for checksum, current supported algorithms: md5, sha1, sha256, sha512, crc32, crc64, adler32, fnv-1-32, fnv-1a-32, fnv-1-64, fnv-1a-64, fnv-1-128, fnv-1a-128")
	cl.BoolVar(&config.Progress, "progress", false, "print the sync progress")
	cl.Int64Var(&config.MaxTranRate, "max_tran_rate", 0, "limit the max transmission rate in the server and client sides, and this is an expected value, not an absolute one")
	cl.BoolVar(&config.DryRun, "dry_run", false, "In dry run mode, gofs is started without actual sync operations")
	cl.BoolVar(&config.CopyLink, "copy_link", false, "transform symlink into referent file, and symlinks that point outside the source tree will be ignored")
	cl.BoolVar(&config.CopyUnsafeLink, "copy_unsafe_link", false, "force to transform the symlinks that point outside the source tree into referent file")

	// ssh
	cl.StringVar(&config.SSHKey, "ssh_key", "", "a cryptographic key used for authenticating computers in the SSH protocol")

	// file monitor
	cl.BoolVar(&config.EnableSyncDelay, "sync_delay", false, "enable sync delay, start sync when the event count is equal or greater than -sync_delay_events, or wait for -sync_delay_time interval time since the last sync")
	cl.IntVar(&config.SyncDelayEvents, "sync_delay_events", 10, "the maximum event count of sync delay")
	cl.DurationVar(&config.SyncDelayTime, "sync_delay_time", time.Second*30, "the maximum delay interval time after the last sync")
	cl.IntVar(&config.SyncWorkers, "sync_workers", 1, "the number of file sync workers")

	// retry
	cl.IntVar(&config.RetryCount, "retry_count", 15, "if execute failed, then retry to work -retry_count times")
	cl.DurationVar(&config.RetryWait, "retry_wait", time.Second*5, "if retry to work, wait -retry_wait time then do")
	cl.BoolVar(&config.RetryAsync, "retry_async", false, "execute retry asynchronously")

	// log
	cl.IntVar(&config.LogLevel, "log_level", int(level.InfoLevel), "set log level, default is INFO. DEBUG=0 INFO=1 WARN=2 ERROR=3")
	cl.BoolVar(&config.EnableFileLogger, "log_file", true, "enable the file logger")
	cl.StringVar(&config.LogDir, "log_dir", "./logs/", "set the directory of the log file")
	cl.BoolVar(&config.LogFlush, "log_flush", true, "enable auto flush log with interval")
	cl.DurationVar(&config.LogFlushInterval, "log_flush_interval", time.Second*3, "set the log flush interval duration, you need to enable -log_flush first")
	cl.BoolVar(&config.EnableEventLog, "log_event", false, "enable the event log")
	cl.Float64Var(&config.LogSampleRate, "log_sample_rate", 1, "set the sample rate for the sample logger, and the value ranges from 0 to 1")
	cl.StringVar(&config.LogFormat, "log_format", formatter.TextFormatter, "set the log output format, current support text and json")
	cl.BoolVar(&config.LogSplitDate, "log_split_date", false, "split log file by date")

	// daemon
	cl.BoolVar(&config.IsDaemon, "daemon", false, "enable daemon to create and monitor a subprocess to work, you can use [go build -ldflags=\"-H windowsgui\"] to build on Windows")
	cl.BoolVar(&config.DaemonPid, "daemon_pid", false, "record parent process pid, daemon process pid and worker process pid to pid file")
	cl.DurationVar(&config.DaemonDelay, "daemon_delay", time.Second, "daemon work interval, wait to create subprocess")
	cl.DurationVar(&config.DaemonMonitorDelay, "daemon_monitor_delay", time.Second*3, "daemon monitor work interval, wait to check subprocess state")
	cl.BoolVar(&config.KillPPid, "kill_ppid", false, "try to kill the parent process when it's running")
	cl.BoolVar(&config.IsSubprocess, daemon.SubprocessTag, false, "tag current process is subprocess")

	// file server
	cl.BoolVar(&config.EnableFileServer, "server", false, "start a file server to browse source directory and dest directory")
	cl.StringVar(&config.FileServerAddr, "server_addr", server.DefaultAddrHttps, "a file server binding address")
	cl.BoolVar(&config.EnableFileServerCompress, "server_compress", false, "enable response compression for the file server")
	cl.BoolVar(&config.EnableManage, "manage", false, "enable the manage api route")
	cl.BoolVar(&config.ManagePrivate, "manage_private", true, "allow to access manage api route by private address and loopback address only")
	cl.BoolVar(&config.EnablePushServer, "push_server", false, "whether to enable the push server")
	cl.BoolVar(&config.EnableReport, "report", false, "enable the report api route and start to collect the report data, need to enable -manage flag first")
	cl.StringVar(&config.SessionConnection, "session_connection", "memory:", "the session connection string, an example for redis session: redis://127.0.0.1:6379?password=redis_password&db=10&max_idle=10&secret=redis_secret")

	// http protocol
	cl.BoolVar(&config.EnableHTTP3, "http3", false, "enable the HTTP3 protocol, pay attention to what you enable the TLS first")

	// tls transfer
	cl.BoolVar(&config.EnableTLS, "tls", true, fmt.Sprintf("enable the tls connections, if disable it, server_addr is \"%s\" default", server.DefaultAddrHttp))
	cl.StringVar(&config.TLSCertFile, "tls_cert_file", "gofs.pem", "cert file for tls connections")
	cl.StringVar(&config.TLSKeyFile, "tls_key_file", "gofs.key", "key file for tls connections")
	cl.BoolVar(&config.TLSInsecureSkipVerify, "tls_insecure_skip_verify", true, "controls whether a client skip verifies the server's certificate chain and host name")

	// login user
	cl.StringVar(&config.Users, "users", "", "the server accounts, the server allows anonymous access if there is no effective account, format like this, user1|password1|rwx,user2|password2|rwx")
	cl.IntVar(&config.RandomUserCount, "rand_user_count", 0, "the number of random server accounts, if it is greater than zero, random generate some accounts for -users")
	cl.IntVar(&config.RandomUserNameLen, "rand_user_len", 6, "the length of the random user's username")
	cl.IntVar(&config.RandomPasswordLen, "rand_pwd_len", 10, "the length of the random user's password")
	cl.StringVar(&config.RandomDefaultPerm, "rand_perm", "r", "the default permission of every random user, like 'rwx'")
	cl.StringVar(&config.TokenSecret, "token_secret", "", "a secret string for token")

	// checksum
	cl.BoolVar(&config.Checksum, "checksum", false, "calculate and print the checksum for source file")

	// encrypt
	cl.BoolVar(&config.Encrypt, "encrypt", false, "enable the encrypt path")
	cl.StringVar(&config.EncryptPath, "encrypt_path", "", "the files in the encrypt path will be encrypted before sync to destination")
	cl.StringVar(&config.EncryptSecret, "encrypt_secret", "", "a secret string for encryption")

	// decrypt
	cl.BoolVar(&config.Decrypt, "decrypt", false, "decrypt the files from decrypt path to decrypt output path")
	cl.StringVar(&config.DecryptPath, "decrypt_path", "", "a directory or file to decrypt")
	cl.StringVar(&config.DecryptSecret, "decrypt_secret", "", "a secret string for decryption")
	cl.StringVar(&config.DecryptOut, "decrypt_out", "", "the decrypt files output directory path")

	// task
	cl.StringVar(&config.TaskConf, "task_conf", "", "the task conf address")
	cl.BoolVar(&config.EnableTaskClient, "task_client", false, "start a task client")
	cl.StringVar(&config.TaskClientLabels, "task_client_labels", "", "the labels of the task client")
	cl.IntVar(&config.TaskClientMaxWorker, "task_client_max_worker", 1, "limit the max concurrent workers in the task client side")

	cl.Parse(args[1:])
	return config
}
