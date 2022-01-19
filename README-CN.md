# gofs

[English](README.md) | 简体中文

基于golang开发的一款开箱即用的文件同步工具

## 安装

```bash
go install github.com/no-src/gofs/...@latest
```

### 后台运行

在windows系统中，你可以使用下面的命令构建一个在后台运行的不带命令行界面的程序

```bat
go install -ldflags="-H windowsgui" github.com/no-src/gofs/...@latest
```

### 移除Web文件服务器

如果你不需要一个Web文件服务器，可以使用下面命令构建一个体积更小的不带Web文件服务器的程序

```bash
go install -tags "no_server" github.com/no-src/gofs/...@latest
```

## 快速开始

### 先决条件

请确保文件同步的源目录和目标目录都已经存在，如果目录不存在，则用你实际的目录替换下面的路径进行提前创建

```bash
$ mkdir src target
```

生成仅用于测试的证书和密钥文件，生产中请替换为正式的证书

TLS证书和密钥文件仅用于与[Web文件服务器](#web文件服务器)和[远程磁盘服务端](#远程磁盘服务端)进行安全通讯

```bash
$ go run $GOROOT/src/crypto/tls/generate_cert.go --host 127.0.0.1
2021/12/30 17:21:54 wrote cert.pem
2021/12/30 17:21:54 wrote key.pem
```

查看你的工作目录

```bash
$ ls
cert.pem  key.pem  src  target
```

### 本地磁盘

监控本地源目录将变更同步到目标目录

你可以使用`logically_delete`命令行参数来启用逻辑删除，从而避免误删数据

```bash
$ gofs -src=./src -target=./target
```

### 全量同步

执行一次全量同步，直接将整个源目录同步到目标目录

```bash
$ gofs -src=./src -target=./target -sync_once
```

### 定时同步

定时执行全量同步，将整个源目录同步到目标目录

```bash
# 每30秒钟将源目录全量同步到目标目录
$ gofs -src=./src -target=./target -sync_cron="*/30 * * * * *"
```

### 守护进程模式

启动守护进程来创建一个工作进程处理实际的任务，并将相关进程的pid信息记录到pid文件中

```bash
$  gofs -src=./src -target=./target -daemon -daemon_pid
```

### Web文件服务器

启动一个Web文件服务器用于访问远程的源目录和目标目录

Web文件服务器默认使用HTTPS协议，使用`tls_cert_file`和`tls_key_file`命令行参数来指定相关的证书和密钥文件

如果你不需要使用TLS进行安全通讯，可以通过将`tls`命令行参数指定为`false`来禁用它

出于安全考虑，你应该设置`rand_user_count`命令行参数来随机生成指定数量的用户或者通过`users`命令行参数自定义用户信息来保证数据的访问安全，禁止用户匿名访问数据

如果`rand_user_count`命令行参数设置大于0，则随机生成的账户密码将会打印到日志信息中，请注意查看

```bash
# 启动一个Web文件服务器并随机创建3个用户
# 在生产环境中请将`tls_cert_file`和`tls_key_file`命令行参数替换为正式的证书和密钥文件
$ gofs -src=./src -target=./target -server -tls_cert_file=cert.pem -tls_key_file=key.pem -rand_user_count=3
```

### 远程磁盘服务端

启动一个远程磁盘服务端作为一个远程文件数据源

```bash
# 启动一个远程磁盘服务端
# 在生产环境中请将`tls_cert_file`和`tls_key_file`命令行参数替换为正式的证书和密钥文件
# 为了安全起见，请使用复杂的账户密码来设置`users`命令行参数
$ gofs -src="rs://127.0.0.1:9016?mode=server&local_sync_disabled=true&path=./src&fs_server=https://127.0.0.1" -target=./target -users="gofs|password" -tls_cert_file=cert.pem -tls_key_file=key.pem
```

### 远程磁盘客户端

启动一个远程磁盘客户端将远程磁盘服务端的文件变更同步到本地目标目录

使用`sync_once`命令行参数，可以直接将远程磁盘服务端的文件整个全量同步到本地目标目录，就跟[全量同步](#全量同步)一样

使用`sync_cron`命令行参数，可以定时将远程磁盘服务端的文件整个全量同步到本地目标目录，就跟[定时同步](#定时同步)一样

```bash
# 启动一个远程磁盘客户端
# 请将`users`命令行参数替换为上面设置的实际账户名密码
$ gofs -src="rs://127.0.0.1:9016" -target=./target -users="gofs|password"
```

### 性能分析

基于[Web文件服务器](#web文件服务器)启用性能分析工具pprof

默认情况下，仅允许私有地址和回环地址访问pprof相关路由

你可以通过将`pprof_private`设置为`false`来禁用默认行为，允许公网IP访问pprof路由

```bash
$ gofs -src=./src -target=./target -server -tls_cert_file=cert.pem -tls_key_file=key.pem -rand_user_count=3 -pprof
```

pprof访问地址如下：

```
https://127.0.0.1/debug/pprof/
```

### 日志

默认情况下会启用文件日志与控制台日志，你可以将`log_file`命令行参数设置为`false`来禁用文件日志

使用`log_level`命令行参数设置日志的等级，默认级别是`INFO`，可选项为：`DEBUG=0` `INFO=1` `WARN=2` `ERROR=3`

使用`log_dir`命令行参数来设置日志文件目录，默认为`./logs/`

使用`log_flush`命令行参数来设置自动刷新日志到文件中，默认启用

使用`log_flush_interval`命令行参数设置自动刷新日志到文件中的频率，默认为`3s`

使用`log_event`命令行参数启用事件日志，所有事件都会记录到文件中，默认为禁用

```bash
# 在"本地磁盘"模式下设置日志信息
$ gofs -src=./src -target=./target -log_file -log_level=0 -log_dir="./logs/" -log_flush -log_flush_interval=3s -log_event
```

## 更多信息

### 帮助信息

```bash
$ gofs -h
```

### 版本信息

```bash
$ gofs -v
```