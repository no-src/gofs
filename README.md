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

Start a remote disk server as a remote file source.

```bash
gofs -src="rs://127.0.0.1:9016?mode=server&path=./src&fs_server=http://127.0.0.1:9015" -target=./target -server
```

Start a remote disk client to sync files from remote disk server.

```bash
gofs -src="rs://127.0.0.1:9016?msg_queue=500" -target=./target
```