# gofs

## Installation

```bash
go install github.com/no-src/gofs/...@latest
```

You can install a no windows gui program using the following command on Windows.

```bat
go install -ldflags="-H windowsgui" github.com/no-src/gofs/...@latest
```

If you needn't file server, you can install the program without the file server to reduce the file size of the program.

```bash
go install -tags "no_server" github.com/no-src/gofs/...@latest
```

## Quick Start

For example, sync src directory to target directory.

```bash
gofs -src=./src -target=./target
```

Start a daemon to create subprocess to work, and record pid info to pid file.

```bash
gofs -src=./src -target=./target -daemon -daemon_pid
```

Start a file server for src path and target path.
The file server is use HTTPS default, set the `tls_cert_file` and `tls_key_file` flags to customize the cert file and key file.
You can disable the HTTPS by set the `server_tls` flag to `false` if you don't need it.

```bash
gofs -src=./src -target=./target -server
```

Start a remote disk server as a remote file source.

```bash
gofs -src="rs://127.0.0.1:9016?mode=server&local_sync_disabled=true&path=./src&fs_server=https://127.0.0.1" -target=./target -server
```

Start a remote disk client to sync files from remote disk server.

```bash
gofs -src="rs://127.0.0.1:9016?msg_queue=500" -target=./target
```