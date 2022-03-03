# gofs

[![Chat](https://img.shields.io/discord/936876326722363472)](https://discord.gg/Ww6bXCcsgv)
[![Build](https://img.shields.io/github/workflow/status/no-src/gofs/Go)](https://github.com/no-src/gofs/actions)
[![License](https://img.shields.io/github/license/no-src/gofs)](https://github.com/no-src/gofs/blob/main/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/no-src/gofs.svg)](https://pkg.go.dev/github.com/no-src/gofs)
[![Go Report Card](https://goreportcard.com/badge/github.com/no-src/gofs)](https://goreportcard.com/report/github.com/no-src/gofs)
[![codecov](https://codecov.io/gh/no-src/gofs/branch/main/graph/badge.svg?token=U5K9HV78P0)](https://codecov.io/gh/no-src/gofs)
[![Release](https://img.shields.io/github/v/release/no-src/gofs)](https://github.com/no-src/gofs/releases)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

English | [简体中文](README-CN.md)

A file synchronization tool out of the box based on golang.

## Installation

```bash
go install github.com/no-src/gofs/...@latest
```

### Run In Docker

If you want to run in a docker, you should install or build with the `-tags netgo` flag or set the environment `CGO_ENABLED=0`, otherwise you may get an error that the `gofs` not found, when the docker container is running.

```bash
go install -tags netgo github.com/no-src/gofs/...@latest
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

Please ensure the source directory and dest directory exists first, replace the following path with your real path.

```bash
$ mkdir source dest
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
cert.pem  key.pem  source  dest
```

### Local Disk

Monitor source directory and sync change files to dest directory.

You can use the `logically_delete` flag to enable the logically delete and avoid deleting files by mistake.

```bash
$ gofs -source=./source -dest=./dest
```

### Sync Once

Sync the whole path immediately from source directory to dest directory.

```bash
$ gofs -source=./source -dest=./dest -sync_once
```

### Sync Cron

Sync the whole path from source directory to dest directory with cron.

```bash
# Per 30 seconds sync the whole path from source directory to dest directory
$ gofs -source=./source -dest=./dest -sync_cron="*/30 * * * * *"
```

### Daemon Mode

Start a daemon to create subprocess to work, and record pid info to pid file.

```bash
$  gofs -source=./source -dest=./dest -daemon -daemon_pid
```

### File Server

Start a file server for source directory and dest directory.

The file server is use HTTPS default, set the `tls_cert_file` and `tls_key_file` flags to customize the cert file and key file.

You can disable the HTTPS by set the `tls` flag to `false` if you don't need it.

If you set the `tls` to `true`, the file server default port is `443`, otherwise it is `80`, and you can customize the default port with the `server_addr` flag, like `-server_addr=":443"`.

You should set the `rand_user_count` flag to auto generate some random users or set the `users` flag to customize server users for security reasons.

The server users will output to log if you set the `rand_user_count` flag greater than zero.

If you need to compress the files, add the `server_compress` flag to enable gzip compression for response, but it is not fast now.

```bash
# Start a file server and create three random users
# Replace the `tls_cert_file` and `tls_key_file` flags with your real cert files in the production environment
$ gofs -source=./source -dest=./dest -server -tls_cert_file=cert.pem -tls_key_file=key.pem -rand_user_count=3
```

### Remote Disk Server

Start a remote disk server as a remote file source.

The `source` flag detail see [Remote Server Source Protocol](#remote-server-source-protocol).

Pay attention to that remote disk server users must have read permission at least, for example, `-users="gofs|password|r"`.

```bash
# Start a remote disk server
# Replace the `tls_cert_file` and `tls_key_file` flags with your real cert files in the production environment
# Replace the `users` flag with complex username and password for security
$ gofs -source="rs://127.0.0.1:8105?mode=server&local_sync_disabled=true&path=./source&fs_server=https://127.0.0.1" -dest=./dest -users="gofs|password|r" -tls_cert_file=cert.pem -tls_key_file=key.pem
```

### Remote Disk Client

Start a remote disk client to sync change files from remote disk server.

Use the `sync_once` flag to sync the whole path immediately from remote disk server to local dest directory, like [Sync Once](#sync-once).

Use the `sync_cron` flag to sync the whole path from remote disk server to local dest directory with cron, like [Sync Cron](#sync-cron).

The `source` flag detail see [Remote Server Source Protocol](#remote-server-source-protocol).

```bash
# Start a remote disk client
# Replace the `users` flag with your real username and password
$ gofs -source="rs://127.0.0.1:8105" -dest=./dest -users="gofs|password"
```

### Remote Push Server

Start a [Remote Disk Server](#remote-disk-server) as a remote file source, then enable the remote push server with the `push_server` flag.

Pay attention to that remote push server users must have read and write permission at least, for example, `-users="gofs|password|rw"`.

```bash
# Start a remote disk server and enable the remote push server
# Replace the `tls_cert_file` and `tls_key_file` flags with your real cert files in the production environment
# Replace the `users` flag with complex username and password for security
$ gofs -source="rs://127.0.0.1:8105?mode=server&local_sync_disabled=true&path=./source&fs_server=https://127.0.0.1" -dest=./dest -users="gofs|password|rw" -tls_cert_file=cert.pem -tls_key_file=key.pem -push_server
```

### Remote Push Client

Start a remote push client to sync change files to the [Remote Push Server](#remote-push-server).

Use the `chunk_size` flag to set the chunk size of the big file to upload. The default value of `chunk_size` is `1048576`, which means `1MB`.

More flag usage see [Remote Disk Client](#remote-disk-client).

```bash
# Start a remote push client and enable local disk sync, sync the file changes from source path to the local dest path and the remote push server
# Replace the `users` flag with your real username and password
$ gofs -source="./source" -dest="rs://127.0.0.1:8105?local_sync_disabled=false&path=./dest" -users="gofs|password"
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

- `path` the [Remote Disk Server](#remote-disk-server) actual local source directory
- `mode` running mode, in [Remote Disk Server](#remote-disk-server) mode is `server`, default is running in [Remote Disk Client](#remote-disk-client) mode
- `fs_server` [File Server](#file-server) address, like `https://127.0.0.1`
- `local_sync_disabled` disabled [Remote Disk Server](#remote-disk-server) sync changes to its local dest path, `true` or `false`, default is `false`

#### Example

For example, in [Remote Disk Server](#remote-disk-server) mode.

```text
 rs://127.0.0.1:8105?mode=server&local_sync_disabled=true&path=./source&fs_server=https://127.0.0.1
 \_/  \_______/ \__/ \____________________________________________________________________________/
  |       |       |                                      |
scheme   host    port                                parameter
```

### Profiling

Enable pprof base [File Server](#file-server).

By default, allow to access pprof route by private address and loopback address only.

You can disable it by setting the `pprof_private` to `false`.

```bash
$ gofs -source=./source -dest=./dest -server -tls_cert_file=cert.pem -tls_key_file=key.pem -rand_user_count=3 -pprof
```

The pprof url address like this

```text
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
$ gofs -source=./source -dest=./dest -log_file -log_level=0 -log_dir="./logs/" -log_flush -log_flush_interval=3s -log_event
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

### About Info

```bash
$ gofs -about
```