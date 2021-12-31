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

```bash
$ gofs -src=./src -target=./target
```

### SyncOnce

Sync the whole path immediately from src directory to target directory.

```bash
$ gofs -src=./src -target=./target -sync_once
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

```bash
# Start a remote disk server
# Replace the `tls_cert_file` and `tls_key_file` flags with your real cert files in the production environment
# Replace the `users` flag with complex username and password for security
$ gofs -src="rs://127.0.0.1:9016?mode=server&local_sync_disabled=true&path=./src&fs_server=https://127.0.0.1" -target=./target -users="gofs|password" -tls_cert_file=cert.pem -tls_key_file=key.pem
```

### Remote Disk Client

Start a remote disk client to sync change files from remote disk server.

You can sync the whole path immediately from remote disk server to local target directory with the `sync_once` flag, like [SyncOnce](#synconce).

```bash
# Start a remote disk client
# Replace the `users` flag with your real username and password
$ gofs -src="rs://127.0.0.1:9016" -target=./target -users="gofs|password"
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