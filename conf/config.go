package conf

import (
	"bytes"
	"os"
	"strings"

	"github.com/no-src/gofs/core"
	"github.com/no-src/nsgo/yamlutil"
)

// Config store all the flag info
type Config struct {
	// other
	PrintVersion bool   `json:"-" yaml:"-"`
	PrintAbout   bool   `json:"-" yaml:"-"`
	Conf         string `json:"-" yaml:"-"`

	// file sync
	Source                core.VFS  `json:"source" yaml:"source"`
	Dest                  core.VFS  `json:"dest" yaml:"dest"`
	SyncOnce              bool      `json:"sync_once" yaml:"sync_once"`
	SyncCron              string    `json:"sync_cron" yaml:"sync_cron"`
	EnableLogicallyDelete bool      `json:"logically_delete" yaml:"logically_delete"`
	ClearDeletedPath      bool      `json:"clear_deleted" yaml:"clear_deleted"`
	IgnoreConf            string    `json:"ignore_conf" yaml:"ignore_conf"`
	IgnoreDeletedPath     bool      `json:"ignore_deleted" yaml:"ignore_deleted"`
	ChunkSize             core.Size `json:"chunk_size" yaml:"chunk_size"`
	CheckpointCount       int       `json:"checkpoint_count" yaml:"checkpoint_count"`
	ForceChecksum         bool      `json:"force_checksum" yaml:"force_checksum"`
	ChecksumAlgorithm     string    `json:"checksum_algorithm" yaml:"checksum_algorithm"`
	Progress              bool      `json:"progress" yaml:"progress"`
	MaxTranRate           core.Size `json:"max_tran_rate" yaml:"max_tran_rate"`
	DryRun                bool      `json:"dry_run" yaml:"dry_run"`
	CopyLink              bool      `json:"copy_link" yaml:"copy_link"`
	CopyUnsafeLink        bool      `json:"copy_unsafe_link" yaml:"copy_unsafe_link"`

	// file monitor
	EnableSyncDelay bool          `json:"sync_delay" yaml:"sync_delay"`
	SyncDelayEvents int           `json:"sync_delay_events" yaml:"sync_delay_events"`
	SyncDelayTime   core.Duration `json:"sync_delay_time" yaml:"sync_delay_time"`
	SyncWorkers     int           `json:"sync_workers" yaml:"sync_workers"`

	// retry
	RetryCount int           `json:"retry_count" yaml:"retry_count"`
	RetryWait  core.Duration `json:"retry_wait" yaml:"retry_wait"`
	RetryAsync bool          `json:"retry_async" yaml:"retry_async"`

	// log
	LogLevel         int           `json:"log_level" yaml:"log_level"`
	EnableFileLogger bool          `json:"log_file" yaml:"log_file"`
	LogDir           string        `json:"log_dir" yaml:"log_dir"`
	LogFlush         bool          `json:"log_flush" yaml:"log_flush"`
	LogFlushInterval core.Duration `json:"log_flush_interval" yaml:"log_flush_interval"`
	EnableEventLog   bool          `json:"log_event" yaml:"log_event"`
	LogSampleRate    float64       `json:"log_sample_rate" yaml:"log_sample_rate"`
	LogFormat        string        `json:"log_format" yaml:"log_format"`
	LogSplitDate     bool          `json:"log_split_date" yaml:"log_split_date"`

	// daemon
	IsDaemon           bool          `json:"daemon" yaml:"daemon"`
	DaemonPid          bool          `json:"daemon_pid" yaml:"daemon_pid"`
	DaemonDelay        core.Duration `json:"daemon_delay" yaml:"daemon_delay"`
	DaemonMonitorDelay core.Duration `json:"daemon_monitor_delay" yaml:"daemon_monitor_delay"`
	KillPPid           bool          `json:"kill_ppid" yaml:"kill_ppid"`
	IsSubprocess       bool          `json:"sub" yaml:"sub"`

	// file server
	EnableFileServer         bool   `json:"server" yaml:"server"`
	FileServerAddr           string `json:"server_addr" yaml:"server_addr"`
	EnableFileServerCompress bool   `json:"server_compress" yaml:"server_compress"`
	EnableManage             bool   `json:"manage" yaml:"manage"`
	ManagePrivate            bool   `json:"manage_private" yaml:"manage_private"`
	EnablePushServer         bool   `json:"push_server" yaml:"push_server"`
	EnableReport             bool   `json:"report" yaml:"report"`
	SessionConnection        string `json:"session_connection" yaml:"session_connection"`

	// http protocol
	EnableHTTP3 bool `json:"http3" yaml:"http3"`

	// tls transfer
	EnableTLS             bool   `json:"tls" yaml:"tls"`
	TLSCertFile           string `json:"tls_cert_file" yaml:"tls_cert_file"`
	TLSKeyFile            string `json:"tls_key_file" yaml:"tls_key_file"`
	TLSInsecureSkipVerify bool   `json:"tls_insecure_skip_verify" yaml:"tls_insecure_skip_verify"`

	// login user
	Users             string `json:"users" yaml:"users"`
	RandomUserCount   int    `json:"rand_user_count" yaml:"rand_user_count"`
	RandomUserNameLen int    `json:"rand_user_len" yaml:"rand_user_len"`
	RandomPasswordLen int    `json:"rand_pwd_len" yaml:"rand_pwd_len"`
	RandomDefaultPerm string `json:"rand_perm" yaml:"rand_perm"`
	TokenSecret       string `json:"token_secret" yaml:"token_secret"`

	// checksum
	Checksum bool `json:"checksum" yaml:"checksum"`

	// encrypt
	Encrypt       bool   `json:"encrypt" yaml:"encrypt"`
	EncryptPath   string `json:"encrypt_path" yaml:"encrypt_path"`
	EncryptSecret string `json:"encrypt_secret" yaml:"encrypt_secret"`

	// decrypt
	Decrypt       bool   `json:"decrypt" yaml:"decrypt"`
	DecryptPath   string `json:"decrypt_path" yaml:"decrypt_path"`
	DecryptSecret string `json:"decrypt_secret" yaml:"decrypt_secret"`
	DecryptOut    string `json:"decrypt_out" yaml:"decrypt_out"`

	// task
	TaskConf            string `json:"task_conf" yaml:"task_conf"`
	EnableTaskClient    bool   `json:"task_client" yaml:"task_client"`
	TaskClientLabels    string `json:"task_client_labels" yaml:"task_client_labels"`
	TaskClientMaxWorker int    `json:"task_client_max_worker" yaml:"task_client_max_worker"`
}

// ToArgs parse the Config to program arguments and the first argument is the current program name
func (c Config) ToArgs() (args []string, err error) {
	data, err := yamlutil.Marshal(c)
	if err != nil {
		return nil, err
	}
	exeFile, err := os.Executable()
	if err == nil {
		args = append(args, exeFile)
		buf := bytes.NewBuffer(data)
		for {
			line, readErr := buf.ReadString('\n')
			if readErr != nil {
				break
			}
			kv := strings.SplitN(line, ": ", 2)
			k := strings.TrimSpace(kv[0])
			v := strings.TrimSpace(kv[1])
			v = strings.Trim(v, "'")
			if v == "\"\"" {
				v = ""
			}
			line = "-" + k + "=" + v
			args = append(args, line)
		}
	}
	return args, err
}
