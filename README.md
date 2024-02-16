# gofs

[![Build](https://img.shields.io/github/actions/workflow/status/no-src/gofs/go.yml?branch=main)](https://github.com/no-src/gofs/actions)
[![License](https://img.shields.io/github/license/no-src/gofs)](https://github.com/no-src/gofs/blob/main/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/no-src/gofs.svg)](https://pkg.go.dev/github.com/no-src/gofs)
[![Go Report Card](https://goreportcard.com/badge/github.com/no-src/gofs)](https://goreportcard.com/report/github.com/no-src/gofs)
[![codecov](https://codecov.io/gh/no-src/gofs/branch/main/graph/badge.svg?token=U5K9HV78P0)](https://codecov.io/gh/no-src/gofs)
[![Release](https://img.shields.io/github/v/release/no-src/gofs)](https://github.com/no-src/gofs/releases)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

English | [简体中文](README-CN.md)

A cross-platform real-time file synchronization tool out of the box based on Golang.

## Installation

The first need [Go](https://go.dev/doc/install) installed (**version 1.21+ is required**), then you can use the below
command to install `gofs`.

```bash
go install github.com/no-src/gofs/...@latest
```

### Run In Docker

You can use the [build-docker.sh](/scripts/build-docker.sh) script to build the docker image and you should clone this
repository and `cd` to the root path of the repository first.

```bash
$ ./scripts/build-docker.sh
```

Or pull the docker image directly from [DockerHub](https://hub.docker.com/r/nosrc/gofs) with the command below.

```bash
$ docker pull nosrc/gofs
```

For more scripts about release and docker, see the [scripts](/scripts) directory.

### Run In the Background

You can install a program run in the background using the following command on Windows.

```bash
go install -ldflags="-H windowsgui" github.com/no-src/gofs/...@latest
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

### Usage

#### Between Disks

Synchronize files between disks by [Local Disk](#local-disk).

```mermaid
sequenceDiagram
    participant DA as DiskA
    participant C as Client
    participant DB as DiskB

    autonumber

    C ->> DA: monitor disk
    DA ->> C: notify change
    C ->> DA: read file
    DA ->> C: return file
    C ->> DB: write file
```

#### From Server

Synchronize files from server by [Remote Disk Server](#remote-disk-server)
and [Remote Disk Client](#remote-disk-client).

```mermaid
sequenceDiagram
    participant SD as Server Disk
    participant S as Server
    participant C as Client
    participant CD as Client Disk

    autonumber

    S ->> SD: monitor disk
    C ->> S: connect and auth
    SD ->> S: notify change
    S ->> C: notify change
    C ->> S: pull file
    S ->> SD: read file
    SD ->> S: return file
    S ->> C: send file
    C ->> CD: write file
```

#### To Server

Synchronize files to server by [Remote Push Server](#remote-push-server) and [Remote Push Client](#remote-push-client).

```mermaid
sequenceDiagram
    participant CD as Client Disk
    participant C as Client
    participant S as Server
    participant SD as Server Disk

    autonumber

    C ->> CD: monitor disk
    CD ->> C: notify change
    C ->> CD: read file
    CD ->> C: return file
    C ->> S: push file
    S ->> SD: write file
```

#### From SFTP Server

Synchronize files from SFTP server by [SFTP Pull Client](#sftp-pull-client).

```mermaid
sequenceDiagram
    participant CD as Client Disk
    participant C as Client
    participant SS as SFTP Server
    participant SSD as SFTP Server Disk

    autonumber

    C ->> SS: pull file
    SS ->> SSD: read file
    SSD ->> SS: return file
    SS ->> C: send file
    C ->> CD: write file
```

#### To SFTP Server

Synchronize files to SFTP server by [SFTP Push Client](#sftp-push-client).

```mermaid
sequenceDiagram
    participant CD as Client Disk
    participant C as Client
    participant SS as SFTP Server
    participant SSD as SFTP Server Disk

    autonumber

    C ->> CD: monitor disk
    CD ->> C: notify change
    C ->> CD: read file
    CD ->> C: return file
    C ->> SS: push file
    SS ->> SSD: write file
```

#### From MinIO Server

Synchronize files from MinIO server by [MinIO Pull Client](#minio-pull-client).

```mermaid
sequenceDiagram
    participant CD as Client Disk
    participant C as Client
    participant MS as MinIO Server
    participant MSD as MinIO Server Disk

    autonumber

    C ->> MS: pull file
    MS ->> MSD: read file
    MSD ->> MS: return file
    MS ->> C: send file
    C ->> CD: write file
```

#### To MinIO Server

Synchronize files to MinIO server by [MinIO Push Client](#minio-push-client).

```mermaid
sequenceDiagram
    participant CD as Client Disk
    participant C as Client
    participant MS as MinIO Server
    participant MSD as MinIO Server Disk

    autonumber

    C ->> CD: monitor disk
    CD ->> C: notify change
    C ->> CD: read file
    CD ->> C: return file
    C ->> MS: push file
    MS ->> MSD: write file
```

#### Task Mode

Start a [Task Client](#task-client) to subscribe to the [Task Server](#task-server), then acquire the task and execute
it, take the [From Server](#from-server) for example.

```mermaid
sequenceDiagram
    participant A as Admin
    participant TS as Task Server
    participant SD as Server Disk
    participant S as Server
    participant TC as Task Client
    participant CW as Client Worker
    participant CD as Client Disk
    participant TQ as Task Queue

    autonumber

    S ->> SD: monitor disk
    S ->> TS: start task server
    A ->> TS: create task
    TC ->> TS: subscribe task
    TS ->> TC: distribute task
    TC ->> CW: start worker
    CW ->> TQ: add to task queue
    TQ ->> CW: execute task
    activate CW
    CW ->> S: connect and auth
    loop
        SD ->> S: notify change
        S ->> CW: notify change
        CW ->> S: pull file
        S ->> SD: read file
        SD ->> S: return file
        S ->> CW: send file
        CW ->> CD: write file
    end
    deactivate CW
```

## Features

### Local Disk

Monitor source directory and sync change files to dest directory.

You can use the `logically_delete` flag to enable the logically delete and avoid deleting files by mistake.

Set the `checkpoint_count` flag to use the checkpoint in the file to reduce transfer unmodified file chunks, by
default `checkpoint_count=10`, which means it has `10+2` checkpoints at most. There are two additional checkpoints at
the head and tail. The first checkpoint is equal to the `chunk_size`, it is optional. The last checkpoint is equal to
the file size, it is required. The checkpoint offset set by the `checkpoint_count` is always more than `chunk_size`,
unless the file size is less than or equal to `chunk_size`, then the `checkpoint_count` will be zero, so it is optional.

By default, if the file size and file modification time of the source file is equal to the destination file, then ignore
the current file transfer. You can use the `force_checksum` flag to force enable the checksum to compare whether the
file is equal or not.

The default checksum hash algorithm is `md5`, you can use the `checksum_algorithm` flag to change the default hash
algorithm, current supported algorithms: `md5`, `sha1`, `sha256`, `sha512`, `crc32`, `crc64`, `adler32`, `fnv-1-32`
, `fnv-1a-32`, `fnv-1-64`, `fnv-1a-64`, `fnv-1-128`, `fnv-1a-128`.

If you want to reduce the frequency of synchronization, you can use the `sync_delay` flag to enable sync delay, start
sync when the event count is equal or greater than `sync_delay_events`, or wait for `sync_delay_time` interval time
since the last sync.

And you can use the `progress` flag to print the file sync progress bar.

```bash
$ gofs -source=./source -dest=./dest
```

### Encryption

You can use `encrypt` flag to enable encryption and specify a directory as an encryption workspace by `encrypt_path`
flag. All files in the directory will be encrypted then sync to the destination path.

```bash
$ gofs -source=./source -dest=./dest -encrypt -encrypt_path=./source/encrypt -encrypt_secret=mysecret_16bytes
```

### Decryption

You can use the `decrypt` flag to decrypt the encryption files to a specified path.

```bash
$ gofs -decrypt -decrypt_path=./dest/encrypt -decrypt_secret=mysecret_16bytes -decrypt_out=./decrypt_out
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
$ gofs -source=./source -dest=./dest -daemon -daemon_pid
```

### File Server

Start a file server for source directory and dest directory.

The file server is use HTTPS default, set the `tls_cert_file` and `tls_key_file` flags to customize the cert file and
key file.

You can disable the HTTPS by set the `tls` flag to `false` if you don't need it.

If you set the `tls` to `true`, the file server default port is `443`, otherwise it is `80`, and you can customize the
default port with the `server_addr` flag, like `-server_addr=":443"`.

If you enable the `tls` flag on the server side, you can control whether a client skip verifies the server's certificate
chain and host name by the `tls_insecure_skip_verify` flag, default is `true`.

If you already enable the `tls` flag, then you can use the `http3` flag to enable the HTTP3 protocol in the server and
client sides.

You should set the `rand_user_count` flag to auto generate some random users or set the `users` flag to customize server
users for security reasons.

The server users will output to log if you set the `rand_user_count` flag greater than zero.

If you need to compress the files, add the `server_compress` flag to enable gzip compression for response, but it is not
fast now, and may reduce transmission efficiency in the LAN.

You can switch the session store mode for the file server by `session_connection` flag,
currently supports memory and redis, default is memory.
If you want to use the redis as the session store, here is an example for redis session connection string:
`redis://127.0.0.1:6379?password=redis_password&db=10&max_idle=10&secret=redis_secret`.

```bash
# Start a file server and create three random users
# Replace the `tls_cert_file` and `tls_key_file` flags with your real cert files in the production environment
$ gofs -source=./source -dest=./dest -server -tls_cert_file=cert.pem -tls_key_file=key.pem -rand_user_count=3
```

### Rate Limit

Use the `max_tran_rate` flag to limit the max transmission rate in the server and client sides,
and this is an expected value, not an absolute one.

For example, limit the max transmission rate to 1048576 bytes, means 1MB.

```bash
$ gofs -source=./source -dest=./dest -max_tran_rate=1048576
```

### Remote Disk Server

Start a remote disk server as a remote file source.

The `source` flag detail see [Remote Server Source Protocol](#remote-server-source-protocol).

Pay attention to that remote disk server users must have read permission at least, for
example, `-users="gofs|password|r"`.

You can use the `checkpoint_count` and `sync_delay` flags like the [Local Disk](#local-disk).

```bash
# Start a remote disk server
# Replace the `tls_cert_file` and `tls_key_file` flags with your real cert files in the production environment
# Replace the `users` flag with complex username and password for security
$ gofs -source="rs://127.0.0.1:8105?mode=server&local_sync_disabled=true&path=./source&fs_server=https://127.0.0.1" -dest=./dest -users="gofs|password|r" -tls_cert_file=cert.pem -tls_key_file=key.pem -token_secret=mysecret_16bytes
```

### Remote Disk Client

Start a remote disk client to sync change files from remote disk server.

The `source` flag detail see [Remote Server Source Protocol](#remote-server-source-protocol).

Use the `sync_once` flag to sync the whole path immediately from remote disk server to local dest directory,
like [Sync Once](#sync-once).

Use the `sync_cron` flag to sync the whole path from remote disk server to local dest directory with cron,
like [Sync Cron](#sync-cron).

Use the `force_checksum` flag to force enable the checksum to compare whether the file is equal or not,
like [Local Disk](#local-disk).

You can use the `sync_delay` flag like the [Local Disk](#local-disk).

```bash
# Start a remote disk client
# Replace the `users` flag with your real username and password
$ gofs -source="rs://127.0.0.1:8105" -dest=./dest -users="gofs|password" -tls_cert_file=cert.pem
```

### Remote Push Server

Start a [Remote Disk Server](#remote-disk-server) as a remote file source, then enable the remote push server with
the `push_server` flag.

Pay attention to that remote push server users must have read and write permission at least, for
example, `-users="gofs|password|rw"`.

```bash
# Start a remote disk server and enable the remote push server
# Replace the `tls_cert_file` and `tls_key_file` flags with your real cert files in the production environment
# Replace the `users` flag with complex username and password for security
$ gofs -source="rs://127.0.0.1:8105?mode=server&local_sync_disabled=true&path=./source&fs_server=https://127.0.0.1" -dest=./dest -users="gofs|password|rw" -tls_cert_file=cert.pem -tls_key_file=key.pem -push_server -token_secret=mysecret_16bytes
```

### Remote Push Client

Start a remote push client to sync change files to the [Remote Push Server](#remote-push-server).

Use the `chunk_size` flag to set the chunk size of the big file to upload. The default value of `chunk_size`
is `1048576`, which means `1MB`.

You can use the `checkpoint_count` and `sync_delay` flags like the [Local Disk](#local-disk).

More flag usage see [Remote Disk Client](#remote-disk-client).

```bash
# Start a remote push client and enable local disk sync, sync the file changes from source path to the local dest path and the remote push server
# Replace the `users` flag with your real username and password
$ gofs -source="./source" -dest="rs://127.0.0.1:8105?local_sync_disabled=false&path=./dest" -users="gofs|password" -tls_cert_file=cert.pem
```

### SFTP Push Client

Start a SFTP push client to sync change files to the SFTP server.

```bash
$ gofs -source="./source" -dest="sftp://127.0.0.1:22?local_sync_disabled=false&path=./dest&remote_path=/gofs_sftp_server&ssh_user=sftp_user&ssh_pass=sftp_pwd"
```

### SFTP Pull Client

Start a SFTP pull client to pull the files from the SFTP server to the local destination path.

```bash
$ gofs -source="sftp://127.0.0.1:22?remote_path=/gofs_sftp_server&ssh_user=sftp_user&ssh_pass=sftp_pwd" -dest="./dest" -sync_once
```

### MinIO Push Client

Start a MinIO push client to sync change files to the MinIO server.

```bash
$ gofs -source="./source" -dest="minio://127.0.0.1:9000?secure=false&local_sync_disabled=false&path=./dest&remote_path=minio-bucket" -users="minio_user|minio_pwd"
```

### MinIO Pull Client

Start a MinIO pull client to pull the files from the MinIO server to the local destination path.

```bash
$ gofs -source="minio://127.0.0.1:9000?secure=false&remote_path=minio-bucket" -dest="./dest" -users="minio_user|minio_pwd" -sync_once
```

### Task Server

Start a task server to distribute the tasks to clients.

Take the [Remote Disk Server](#remote-disk-server) for example, create a task manifest config file like
the [remote-disk-task.yaml](/integration/testdata/conf/task/remote-disk-task.yaml) file first.
Here defined a task that synchronizes files from server.

Then create the task content config
file [run-gofs-remote-disk-client.yaml](/integration/testdata/conf/run-gofs-remote-disk-client.yaml) that defined in the
above manifest config file, and it will be executed by client.

Finally, start the remote disk server with the `task_conf` flag.

Here use `conf` to simplify the command and reuse the integration test config files.

```bash
$ cd integration
$ mkdir -p rs/source rs/dest
$ gofs -conf=./testdata/conf/run-gofs-remote-disk-server.yaml
```

### Task Client

Start a task client to subscribe to the task server, then acquire the task and execute it.

Use the `task_client` flag to start the task client, and the `task_client_max_worker` flag will limit the max
concurrent workers in the task client side.

And you can use the `task_client_labels` flag to define the labels of the task client that use to match the task in the
task server side.

Here use `conf` to simplify the command and reuse the integration test config files.

```bash
$ cd integration
$ mkdir -p rc/source rc/dest
$ gofs -conf=./testdata/conf/run-gofs-task-client.yaml
```

### Relay

If you need to synchronize files between two devices that are unable to establish a direct connection, you can use a
reverse proxy as a relay server. In more detail, see also [Relay](/relay/README.md).

### Remote Server Source Protocol

The remote server source protocol is based on URI, see [RFC 3986](https://www.rfc-editor.org/rfc/rfc3986.html).

#### Scheme

The scheme name is `rs`.

#### Host

The remote server source uses `0.0.0.0` or other local ip address as host in [Remote Disk Server](#remote-disk-server)
mode, and use ip address or domain name as host in [Remote Disk Client](#remote-disk-client) mode.

#### Port

The remote server source port, default is `8105`.

#### Parameter

Use the following parameters in [Remote Disk Server](#remote-disk-server) mode only.

- `path` the [Remote Disk Server](#remote-disk-server) actual local source directory
- `mode` running mode, in [Remote Disk Server](#remote-disk-server) mode is `server`, default is running
  in [Remote Disk Client](#remote-disk-client) mode
- `fs_server` [File Server](#file-server) address, like `https://127.0.0.1`
- `local_sync_disabled` disabled [Remote Disk Server](#remote-disk-server) sync changes to its local dest path, `true`
  or `false`, default is `false`

#### Example

For example, in [Remote Disk Server](#remote-disk-server) mode.

```text
 rs://127.0.0.1:8105?mode=server&local_sync_disabled=true&path=./source&fs_server=https://127.0.0.1
 \_/  \_______/ \__/ \____________________________________________________________________________/
  |       |       |                                      |
scheme   host    port                                parameter
```

### Manage API

Enable manage api base on [File Server](#file-server) by using the `manage` flag.

By default, allow to access manage api by private address and loopback address only.

You can disable it by setting the `manage_private` flag to `false`.

```bash
$ gofs -source=./source -dest=./dest -server -tls_cert_file=cert.pem -tls_key_file=key.pem -rand_user_count=3 -manage
```

#### Profiling API

The pprof url address like this

```text
https://127.0.0.1/manage/pprof/
```

#### Config API

Reading the program config, default return the config with `json` format, and support `json` and `yaml` format
currently.

```text
https://127.0.0.1/manage/config
```

Or use the `format` parameter to specific the config format.

```text
https://127.0.0.1/manage/config?format=yaml
```

#### Report API

Use the `report` flag to enable report api route, and start to collect the report data, need to enable the `manage` flag
first.

The details of the report api see [Report API](/server/README.md#report-api).

```text
https://127.0.0.1/manage/report
```

### Logger

Enable the file logger and console logger by default, and you can disable the file logger by setting the `log_file` flag
to `false`.

Use the `log_level` flag to set the log level, default is `INFO`, (`DEBUG=0` `INFO=1` `WARN=2` `ERROR=3`).

Use the `log_dir` flag to set the directory of the log file, default is `./logs/`.

Use the `log_flush` flag to enable auto flush log with interval, default is `true`.

Use the `log_flush_interval` flag to set the log flush interval duration, default is `3s`.

Use the `log_event` flag to enable the event log, write to file, default is `false`.

Use the `log_sample_rate` flag to set the sample rate for the sample logger, and the value ranges from 0 to 1, default
is `1`.

Use the `log_format` flag to set the log output format, current support `text` and `json`, default is `text`.

Use the `log_split_date` flag to split log file by date, default is `false`.

```bash
# set the logger config in "Local Disk" mode
$ gofs -source=./source -dest=./dest -log_file -log_level=0 -log_dir="./logs/" -log_flush -log_flush_interval=3s -log_event
```

### Use Configuration File

If you want, you can use a configuration file to replace all the flags.It supports `json` and `yaml` format currently.

All the configuration fields are the same as the flags, you can refer to the [Configuration Example](/conf/example)
or the response of [Config API](#config-api).

```bash
$ gofs -conf=./gofs.yaml
```

### Checksum

You can use the `checksum` flag to calculate the file checksum and print the result.

The `chunk_size`, `checkpoint_count` and `checksum_algorithm` flags are effective here the same as in
the [Local Disk](#local-disk).

```bash
$ gofs -source=./gofs -checksum
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

### Web UI

The [gofs-webui](https://github.com/no-src/gofs-webui) is a web UI tool for `gofs`, and it allows you to generate the
config file of `gofs` through the web UI, making the use of gofs easier.

## Contributing

### How To Contribute

```mermaid
%%{init: { "flowchart": {"htmlLabels": false}} }%%
flowchart TD
  PR[pull request]
  MainRepo[github.com/no-src/gofs.git] -- 1.fork --> ForkRepo[github.com/yourname/gofs.git]
  ForkRepo -- 2.git clone --> LocalRepo[local repository]
  LocalRepo -- 3.commit changes --> NewBranch[new branch]
  NewBranch -- 4.git push --> ForkRepo
  ForkRepo -- 5.create pull request --> PR
  PR -- 6.merge to --> MainRepo
```

### Cloud Development Environments

Quick to contribute using cloud development environments.

#### Github Codespaces

[![Open in GitHub Codespaces](https://github.com/codespaces/badge.svg)](https://codespaces.new/no-src/gofs)

#### Gitpod

[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/no-src/gofs)