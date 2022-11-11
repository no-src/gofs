package cmd

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

func parseFlags(args []string) (config conf.Config) {
	if len(args) < 1 {
		panic("at least one argument is required, starting with the program name")
	}

	// print help info if no arguments
	if len(args) <= 1 {
		args = append(args, "-h")
	}

	// other
	flag.BoolVar(&config.PrintVersion, "v", false, "print the version info")
	flag.BoolVar(&config.PrintAbout, "about", false, "print the about info")
	flag.StringVar(&config.Conf, "conf", "", "the path of config file")

	// file sync
	core.VFSVar(&config.Source, "source", core.NewEmptyVFS(), "the source path by monitor")
	core.VFSVar(&config.Dest, "dest", core.NewEmptyVFS(), "the dest path to backup")
	flag.BoolVar(&config.SyncOnce, "sync_once", false, "sync source directory to dest directory once")
	flag.StringVar(&config.SyncCron, "sync_cron", "", "sync source directory to dest directory with cron")
	flag.BoolVar(&config.EnableLogicallyDelete, "logically_delete", false, "delete dest file logically")
	flag.BoolVar(&config.ClearDeletedPath, "clear_deleted", false, "remove all of the deleted files in the dest path")
	flag.StringVar(&config.IgnoreConf, "ignore_conf", "", "a config file of the ignore component")
	flag.BoolVar(&config.IgnoreDeletedPath, "ignore_deleted", true, "ignore to sync the deleted file")
	flag.Int64Var(&config.ChunkSize, "chunk_size", 1024*1024, "the chunk size of the big file")
	flag.IntVar(&config.CheckpointCount, "checkpoint_count", 10, "use the checkpoint in the file to reduce transfer unmodified file chunks")
	flag.BoolVar(&config.ForceChecksum, "force_checksum", false, "if the file size and file modification time of the source file is equal to the destination file and -force_checksum is false, then ignore the current file transfer")
	flag.StringVar(&config.ChecksumAlgorithm, "checksum_algorithm", hashutil.DefaultHash, "set the default hash algorithm for checksum, current supported algorithms: md5, sha1, sha256, sha512, crc32, crc64, adler32, fnv-1-32, fnv-1a-32, fnv-1-64, fnv-1a-64, fnv-1-128, fnv-1a-128")
	flag.BoolVar(&config.Progress, "progress", false, "print the sync progress")

	// ssh
	flag.StringVar(&config.SSHKey, "ssh_key", "", "a cryptographic key used for authenticating computers in the SSH protocol")

	// file monitor
	flag.BoolVar(&config.EnableSyncDelay, "sync_delay", false, "enable sync delay, start sync when the event count is equal or greater than -sync_delay_events, or wait for -sync_delay_time interval time since the last sync")
	flag.IntVar(&config.SyncDelayEvents, "sync_delay_events", 10, "the maximum event count of sync delay")
	core.DurationVar(&config.SyncDelayTime, "sync_delay_time", time.Second*30, "the maximum delay interval time after the last sync")
	flag.IntVar(&config.SyncWorkers, "sync_workers", 1, "the number of file sync workers")

	// retry
	flag.IntVar(&config.RetryCount, "retry_count", 15, "if execute failed, then retry to work -retry_count times")
	core.DurationVar(&config.RetryWait, "retry_wait", time.Second*5, "if retry to work, wait -retry_wait time then do")
	flag.BoolVar(&config.RetryAsync, "retry_async", false, "execute retry asynchronously")

	// log
	flag.IntVar(&config.LogLevel, "log_level", int(level.InfoLevel), "set log level, default is INFO. DEBUG=0 INFO=1 WARN=2 ERROR=3")
	flag.BoolVar(&config.EnableFileLogger, "log_file", true, "enable the file logger")
	flag.StringVar(&config.LogDir, "log_dir", "./logs/", "set the directory of the log file")
	flag.BoolVar(&config.LogFlush, "log_flush", true, "enable auto flush log with interval")
	core.DurationVar(&config.LogFlushInterval, "log_flush_interval", time.Second*3, "set the log flush interval duration, you need to enable -log_flush first")
	flag.BoolVar(&config.EnableEventLog, "log_event", false, "enable the event log")
	flag.Float64Var(&config.LogSampleRate, "log_sample_rate", 1, "set the sample rate for the sample logger, and the value ranges from 0 to 1")
	flag.StringVar(&config.LogFormat, "log_format", formatter.TextFormatter, "set the log output format, current support text and json")
	flag.BoolVar(&config.LogSplitDate, "log_split_date", false, "split log file by date")

	// daemon
	flag.BoolVar(&config.IsDaemon, "daemon", false, "enable daemon to create and monitor a subprocess to work, you can use [go build -ldflags=\"-H windowsgui\"] to build on Windows")
	flag.BoolVar(&config.DaemonPid, "daemon_pid", false, "record parent process pid, daemon process pid and worker process pid to pid file")
	core.DurationVar(&config.DaemonDelay, "daemon_delay", time.Second, "daemon work interval, wait to create subprocess")
	core.DurationVar(&config.DaemonMonitorDelay, "daemon_monitor_delay", time.Second*3, "daemon monitor work interval, wait to check subprocess state")
	flag.BoolVar(&config.KillPPid, "kill_ppid", false, "try to kill the parent process when it's running")
	flag.BoolVar(&config.IsSubprocess, daemon.SubprocessTag, false, "tag current process is subprocess")

	// file server
	flag.BoolVar(&config.EnableFileServer, "server", false, "start a file server to browse source directory and dest directory")
	flag.StringVar(&config.FileServerAddr, "server_addr", server.DefaultAddrHttps, "a file server binding address")
	flag.BoolVar(&config.EnableFileServerCompress, "server_compress", false, "enable response compression for the file server")
	flag.BoolVar(&config.EnableManage, "manage", false, "enable the manage api route")
	flag.BoolVar(&config.ManagePrivate, "manage_private", true, "allow to access manage api route by private address and loopback address only")
	flag.BoolVar(&config.EnablePushServer, "push_server", false, "whether to enable the push server")
	flag.BoolVar(&config.EnableReport, "report", false, "enable the report api route and start to collect the report data, need to enable -manage flag first")
	flag.IntVar(&config.SessionMode, "session_mode", server.MemorySession, "the session store mode for the file server, currently supports memory[1] and redis[2], default is memory[1]")
	flag.StringVar(&config.SessionConnection, "session_connection", "", "the session connection string, an example for redis session: redis://127.0.0.1:6379?password=redis_password&db=10&max_idle=10&secret=redis_secret")

	// tls transfer
	flag.BoolVar(&config.EnableTLS, "tls", true, fmt.Sprintf("enable the tls connections, if disable it, server_addr is \"%s\" default", server.DefaultAddrHttp))
	flag.StringVar(&config.TLSCertFile, "tls_cert_file", "gofs.pem", "cert file for tls connections")
	flag.StringVar(&config.TLSKeyFile, "tls_key_file", "gofs.key", "key file for tls connections")
	flag.BoolVar(&config.TLSInsecureSkipVerify, "tls_insecure_skip_verify", true, "controls whether a client skip verifies the server's certificate chain and host name")

	// login user
	flag.StringVar(&config.Users, "users", "", "the server accounts, the server allows anonymous access if there is no effective account, format like this, user1|password1|rwx,user2|password2|rwx")
	flag.IntVar(&config.RandomUserCount, "rand_user_count", 0, "the number of random server accounts, if it is greater than zero, random generate some accounts for -users")
	flag.IntVar(&config.RandomUserNameLen, "rand_user_len", 6, "the length of the random user's username")
	flag.IntVar(&config.RandomPasswordLen, "rand_pwd_len", 10, "the length of the random user's password")
	flag.StringVar(&config.RandomDefaultPerm, "rand_perm", "r", "the default permission of every random user, like 'rwx'")

	// checksum
	flag.BoolVar(&config.Checksum, "checksum", false, "calculate and print the checksum for source file")

	// encrypt
	flag.BoolVar(&config.Encrypt, "encrypt", false, "enable the encrypt path")
	flag.StringVar(&config.EncryptPath, "encrypt_path", "", "the files in the encrypt path will be encrypted before sync to destination")
	flag.StringVar(&config.EncryptSecret, "encrypt_secret", "", "a secret string for encryption")

	// decrypt
	flag.BoolVar(&config.Decrypt, "decrypt", false, "decrypt the files from decrypt path to decrypt output path")
	flag.StringVar(&config.DecryptPath, "decrypt_path", "", "a directory or file to decrypt")
	flag.StringVar(&config.DecryptSecret, "decrypt_secret", "", "a secret string for decryption")
	flag.StringVar(&config.DecryptOut, "decrypt_out", "", "the decrypt files output directory path")

	flag.CommandLine.Parse(args[1:])

	return config
}
