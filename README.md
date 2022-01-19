# gofs

English | [简体中文](README-CN.md)

A file synchronization tool out of the box based on golang.

## Installation

```bash
go install github.com/no-src/gofs/...@latest
```

### Run In the Background

You can install a program run in the background using the following command on Windows.

```bat
go install -ldflags="-H windowsgui" github.com/no-src/gofs/...@latest
```

### Remove File Server

If you don't need a file server, you can install the program without the file server to reduce the file size of this.

```bash
go install -tags "no_server" github.com/no-src/gofs/...@latest
```

## Quick Start

### Prerequisites

Please ensure the src directory and target directory exists first, replace the following path with your real path.

```bash
$ mkdir src target
```

Generate the TLS cert file and key file for testing purposes.

The TLS cert and key files are just used by [File Server](#file-server) and [Remote Disk Server](#remote-disk-server).

```bash
$ go run $GOROOT/src/crypto/tls/generate_cert.go --host 127.0.0.1
2021/12/30 17:21:54 wrote cert.pem
2021/12/30 17:21:54 wrote key.pem
```

Look up our workspace.

```bash
$ ls
cert.pem  key.pem  src  target
```

### Local Disk

Monitor src directory and sync change files to target directory.

You can use the `logically_delete` flag to enable the logically delete and avoid deleting files by mistake.

```bash
$ gofs -src=./src -target=./target
```

### Sync Once

Sync the whole path immediately from src directory to target directory.

```bash
$ gofs -src=./src -target=./target -sync_once
```

### Sync Cron

Sync the whole path from src directory to target directory with cron.

```bash
# Per 30 seconds sync the whole path from src directory to target directory
$ gofs -src=./src -target=./target -sync_cron="*/30 * * * * *"
```

### Daemon Mode

Start a daemon to create subprocess to work, and record pid info to pid file.

```bash
$  gofs -src=./src -target=./target -daemon -daemon_pid
```

### File Server

Start a file server for src directory and target directory.

The file server is use HTTPS default, set the `tls_cert_file` and `tls_key_file` flags to customize the cert file and key file.

You can disable the HTTPS by set the `tls` flag to `false` if you don't need it.

You should set the `rand_user_count` flag to auto generate some random users or set the `users` flag to customize server users for security reasons.

The server users will output to log if you set the `rand_user_count` flag greater than zero.

```bash
# Start a file server and create three random users
# Replace the `tls_cert_file` and `tls_key_file` flags with your real cert files in the production environment
$ gofs -src=./src -target=./target -server -tls_cert_file=cert.pem -tls_key_file=key.pem -rand_user_count=3
```

### Remote Disk Server

Start a remote disk server as a remote file source.

The `src` flag detail see [Remote Server Source Protocol](#remote-server-source-protocol).

```bash
# Start a remote disk server
# Replace the `tls_cert_file` and `tls_key_file` flags with your real cert files in the production environment
# Replace the `users` flag with complex username and password for security
$ gofs -src="rs://127.0.0.1:8105?mode=server&local_sync_disabled=true&path=./src&fs_server=https://127.0.0.1" -target=./target -users="gofs|password" -tls_cert_file=cert.pem -tls_key_file=key.pem
```

### Remote Disk Client

Start a remote disk client to sync change files from remote disk server.

Use the `sync_once` flag to sync the whole path immediately from remote disk server to local target directory, like [Sync Once](#sync-once).

Use the `sync_cron` flag to sync the whole path from remote disk server to local target directory with cron, like [Sync Cron](#sync-cron).

The `src` flag detail see [Remote Server Source Protocol](#remote-server-source-protocol).

```bash
# Start a remote disk client
# Replace the `users` flag with your real username and password
$ gofs -src="rs://127.0.0.1:8105" -target=./target -users="gofs|password"
```

### Remote Server Source Protocol

The remote server source protocol is based on URI, see [RFC 3986](https://www.rfc-editor.org/rfc/rfc3986.html).

#### Scheme

The scheme name is `rs`.

#### Host

The remote server source uses `0.0.0.0` or other local ip address as host in [Remote Disk Server](#remote-disk-server) mode, and 
use ip address or domain name as host in [Remote Disk Client](#remote-disk-client) mode.

#### Port

The remote server source port, default is `8105`.

#### Parameter

Use the following parameters in [Remote Disk Server](#remote-disk-server) mode only.

- `path` the [Remote Disk Server](#remote-disk-server) actual local src directory
- `mode` running mode, in [Remote Disk Server](#remote-disk-server) mode is `server`, default is running in [Remote Disk Client](#remote-disk-client) mode
- `fs_server` [File Server](#file-server) address, like `https://127.0.0.1`
- `local_sync_disabled` disabled [Remote Disk Server](#remote-disk-server) sync changes to its local target path, `true` or `false`, default is `false`

#### Example

For example, in [Remote Disk Server](#remote-disk-server) mode.

```text
 rs://127.0.0.1:8105?mode=server&local_sync_disabled=true&path=./src&fs_server=https://127.0.0.1
 \_/  \_______/ \__/ \_________________________________________________________________________/
  |       |       |                                      |
scheme   host    port                                parameter
```

### Profiling

Enable pprof base [File Server](#file-server).

By default, allow to access pprof route by private address and loopback address only.

You can disable it by setting the `pprof_private` to `false`.

```bash
$ gofs -src=./src -target=./target -server -tls_cert_file=cert.pem -tls_key_file=key.pem -rand_user_count=3 -pprof
```

The pprof url address like this

```
https://127.0.0.1/debug/pprof/
```

### Logger

Enable the file logger and console logger by default, and you can disable the file logger by setting the `log_file` flag to `false`.

Use the `log_level` flag to set the log level, default is `INFO`, (`DEBUG=0` `INFO=1` `WARN=2` `ERROR=3`).

Use the `log_dir` flag to set the directory of the log file, default is `./logs/`.

Use the `log_flush` flag to enable auto flush log with interval, default is `true`.

Use the `log_flush_interval` flag to set the log flush interval duration, default is `3s`.

Use the `log_event` flag to enable the event log, write to file, default is `false`.

```bash
# set the logger config in "Local Disk" mode
$ gofs -src=./src -target=./target -log_file -log_level=0 -log_dir="./logs/" -log_flush -log_flush_interval=3s -log_event
```

## For More Information

### Help Info

```bash
$ gofs -h
```

### Version Info

```bash
$ gofs -v
```