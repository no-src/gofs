# gofs

## Installation

```bash
go install github.com/no-src/gofs/...@latest
```

You can install a no windows gui program using the following command on Windows.

```bat
go install -ldflags="-H windowsgui" github.com/no-src/gofs/...@latest
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