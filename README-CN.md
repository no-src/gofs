# gofs

[![Chat](https://img.shields.io/discord/936876326722363472)](https://discord.gg/Ww6bXCcsgv)
[![Build](https://img.shields.io/github/workflow/status/no-src/gofs/Go)](https://github.com/no-src/gofs/actions)
[![License](https://img.shields.io/github/license/no-src/gofs)](https://github.com/no-src/gofs/blob/main/LICENSE)
[![Go Reference](https://pkg.go.dev/badge/github.com/no-src/gofs.svg)](https://pkg.go.dev/github.com/no-src/gofs)
[![Go Report Card](https://goreportcard.com/badge/github.com/no-src/gofs)](https://goreportcard.com/report/github.com/no-src/gofs)
[![codecov](https://codecov.io/gh/no-src/gofs/branch/main/graph/badge.svg?token=U5K9HV78P0)](https://codecov.io/gh/no-src/gofs)
[![Release](https://img.shields.io/github/v/release/no-src/gofs)](https://github.com/no-src/gofs/releases)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

[English](README.md) | 简体中文

基于golang开发的一款开箱即用的跨平台文件同步工具

## 安装

首先需要确保已经安装了[Go](https://golang.google.cn/doc/install) (**版本必须是1.18+**)，然后你就可以使用下面的命令来安装`gofs`了

```bash
go install github.com/no-src/gofs/...@latest
```

### 在Docker中运行

如果你想要在docker中运行，请在安装或者编译的时候使用`-tags netgo`命令行参数或者设置环境变量`CGO_ENABLED=0`,否则当你运行docker容器时可能会报找不到`gofs`的错误

```bash
go install -tags netgo github.com/no-src/gofs/...@latest
```

你可以使用[build-docker.sh](/scripts/build-docker.sh)脚本来构建docker镜像，首先你需要克隆本仓库并且`cd`到本仓库的根目录

```bash
$ ./scripts/build-docker.sh
```

或者使用以下命令直接从[DockerHub](https://hub.docker.com/r/nosrc/gofs)中拉取docker镜像

```bash
$ docker pull nosrc/gofs
```

更多关于发布与docker的脚本参见[scripts](scripts)目录

### 后台运行

在windows系统中，你可以使用下面的命令构建一个在后台运行的不带命令行界面的程序

```bat
go install -ldflags="-H windowsgui" github.com/no-src/gofs/...@latest
```

## 快速开始

### 先决条件

请确保文件同步的源目录和目标目录都已经存在，如果目录不存在，则用你实际的目录替换下面的路径进行提前创建

```bash
$ mkdir source dest
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
cert.pem  key.pem  source  dest
```

### 本地磁盘

监控本地源目录将变更同步到目标目录

你可以使用`logically_delete`命令行参数来启用逻辑删除，从而避免误删数据

设置`checkpoint_count`命令行参数来使用文件中的检查点来减少传输未修改的文件块，默认情况下`checkpoint_count=10`
，这意味着它最多有`10+2`个检查点。在头部和尾部还有两个额外的检查点。第一个检查点等于`chunk_size`
，它是可选的。最后一个检查点等于文件大小，这是必需的。由`checkpoint_count`设置的检查点偏移量总是大于`chunk_size`，除非文件大小小于或等于`chunk_size`，那么`checkpoint_count`
将变为`0`，所以它是可选的

默认情况下，如果源文件的大小和修改时间与目标文件相同，则忽略当前文件的传输。你可以使用`force_checksum`命令行参数强制启用校验和来比较文件是否相等

```bash
$ gofs -source=./source -dest=./dest
```

### 全量同步

执行一次全量同步，直接将整个源目录同步到目标目录

```bash
$ gofs -source=./source -dest=./dest -sync_once
```

### 定时同步

定时执行全量同步，将整个源目录同步到目标目录

```bash
# 每30秒钟将源目录全量同步到目标目录
$ gofs -source=./source -dest=./dest -sync_cron="*/30 * * * * *"
```

### 守护进程模式

启动守护进程来创建一个工作进程处理实际的任务，并将相关进程的pid信息记录到pid文件中

```bash
$ gofs -source=./source -dest=./dest -daemon -daemon_pid
```

### Web文件服务器

启动一个Web文件服务器用于访问远程的源目录和目标目录

Web文件服务器默认使用HTTPS协议，使用`tls_cert_file`和`tls_key_file`命令行参数来指定相关的证书和密钥文件

如果你不需要使用TLS进行安全通讯，可以通过将`tls`命令行参数指定为`false`来禁用它

如果将`tls`设置为`true`，则服务器默认运行端口为`443`，反之默认端口为`80`，你可以使用`server_addr`命令行参数来自定义服务器运行端口，例如`-server_addr=":443"`

出于安全考虑，你应该设置`rand_user_count`命令行参数来随机生成指定数量的用户或者通过`users`命令行参数自定义用户信息来保证数据的访问安全，禁止用户匿名访问数据

如果`rand_user_count`命令行参数设置大于0，则随机生成的账户密码将会打印到日志信息中，请注意查看

如果你需要启用gzip压缩响应结果，则添加`server_compress`命令行参数，但是目前gzip压缩不是很快，在局域网中可能会影响传输效率

```bash
# 启动一个Web文件服务器并随机创建3个用户
# 在生产环境中请将`tls_cert_file`和`tls_key_file`命令行参数替换为正式的证书和密钥文件
$ gofs -source=./source -dest=./dest -server -tls_cert_file=cert.pem -tls_key_file=key.pem -rand_user_count=3
```

### 远程磁盘服务端

启动一个远程磁盘服务端作为一个远程文件数据源

`source`命令行参数详见[远程磁盘服务端数据源协议](#远程磁盘服务端数据源协议)

注意远程磁盘服务端的用户至少要拥有读权限，例如：`-users="gofs|password|r"`

你可以使用`checkpoint_count`命令行参数就跟[本地磁盘](#本地磁盘)一样

```bash
# 启动一个远程磁盘服务端
# 在生产环境中请将`tls_cert_file`和`tls_key_file`命令行参数替换为正式的证书和密钥文件
# 为了安全起见，请使用复杂的账户密码来设置`users`命令行参数
$ gofs -source="rs://127.0.0.1:8105?mode=server&local_sync_disabled=true&path=./source&fs_server=https://127.0.0.1" -dest=./dest -users="gofs|password|r" -tls_cert_file=cert.pem -tls_key_file=key.pem
```

### 远程磁盘客户端

启动一个远程磁盘客户端将远程磁盘服务端的文件变更同步到本地目标目录

使用`sync_once`命令行参数，可以直接将远程磁盘服务端的文件整个全量同步到本地目标目录，就跟[全量同步](#全量同步)一样

使用`sync_cron`命令行参数，可以定时将远程磁盘服务端的文件整个全量同步到本地目标目录，就跟[定时同步](#定时同步)一样

使用`force_checksum`命令行参数强制启用校验和来比较文件是否相等，就跟[本地磁盘](#本地磁盘)一样

`source`命令行参数详见[远程磁盘服务端数据源协议](#远程磁盘服务端数据源协议)

```bash
# 启动一个远程磁盘客户端
# 请将`users`命令行参数替换为上面设置的实际账户名密码
$ gofs -source="rs://127.0.0.1:8105" -dest=./dest -users="gofs|password"
```

### 远程推送服务端

启动一个[远程磁盘服务端](#远程磁盘服务端)作为一个远程文件数据源，并使用`push_server`命令行参数启用远程推送服务端

注意远程推送服务端的用户至少要拥有读写权限，例如：`-users="gofs|password|rw"`

```bash
# 启动一个远程磁盘服务端并启用远程推送服务端
# 在生产环境中请将`tls_cert_file`和`tls_key_file`命令行参数替换为正式的证书和密钥文件
# 为了安全起见，请使用复杂的账户密码来设置`users`命令行参数
$ gofs -source="rs://127.0.0.1:8105?mode=server&local_sync_disabled=true&path=./source&fs_server=https://127.0.0.1" -dest=./dest -users="gofs|password|rw" -tls_cert_file=cert.pem -tls_key_file=key.pem -push_server
```

### 远程推送客户端

启动一个远程推送客户端将本地文件变更同步到[远程推送服务端](#远程推送服务端)

使用`chunk_size`命令行参数来设置大文件上传时切分的区块大小，默认值为`1048576`，即`1MB`

你可以使用`checkpoint_count`命令行参数就跟[本地磁盘](#本地磁盘)一样

更多命令行参数用法请参见[远程磁盘客户端](#远程磁盘客户端)

```bash
# 启动一个远程推送客户端并且启用本地磁盘同步，将source目录下的文件变更同步到本地dest目录和远程推送服务器上
# 请将`users`命令行参数替换为上面设置的实际账户名密码
$ gofs -source="./source" -dest="rs://127.0.0.1:8105?local_sync_disabled=false&path=./dest" -users="gofs|password"
```

### 远程磁盘服务端数据源协议

远程磁盘服务端数据源协议基于URI基本语法,详见[RFC 3986](https://www.rfc-editor.org/rfc/rfc3986.html)

#### 方案

方案名称为`rs`

#### 主机名

远程磁盘服务端数据源在[远程磁盘服务端](#远程磁盘服务端)模式下使用`0.0.0.0`或者其他本地网卡IP地址作为主机名，在[远程磁盘客户端](#远程磁盘客户端)
模式下使用远程磁盘服务端的IP地址或者域名作为主机名

#### 端口号

远程磁盘服务端数据源端口号，默认为`8105`

#### 参数

仅在[远程磁盘服务端](#远程磁盘服务端)模式下设置以下参数

- `path` [远程磁盘服务端](#远程磁盘服务端)真实的本地源目录
- `mode` 指定运行模式，只有在[远程磁盘服务端](#远程磁盘服务端)模式下需要手动指定为`server`，默认为[远程磁盘客户端](#远程磁盘客户端)模式
- `fs_server` [Web文件服务器](#web文件服务器)地址，例如`https://127.0.0.1`
- `local_sync_disabled` 是否将[远程磁盘服务端](#远程磁盘服务端)的文件变更同步到远程本地的目标目录,可选值为`true`或`false`，默认值为`false`

#### 示例

[远程磁盘服务端](#远程磁盘服务端)模式下的示例

```text
 rs://127.0.0.1:8105?mode=server&local_sync_disabled=true&path=./source&fs_server=https://127.0.0.1
 \_/  \_______/ \__/ \____________________________________________________________________________/
  |       |       |                                      |
 方案   主机名   端口号                                    参数
```

### 管理接口

基于[Web文件服务器](#web文件服务器)的应用管理接口

默认情况下，仅允许私有地址和回环地址访问管理接口的相关路由

你可以通过将`manage_private`命令行参数设置为`false`来禁用默认行为，允许公网IP访问管理接口的路由

```bash
$ gofs -source=./source -dest=./dest -server -tls_cert_file=cert.pem -tls_key_file=key.pem -rand_user_count=3 -manage
```

#### 性能分析接口

pprof访问地址如下：

```text
https://127.0.0.1/manage/pprof/
```

#### 配置接口

读取应用程序配置，默认返回`json`格式，当前支持`json`和`yaml`格式

```text
https://127.0.0.1/manage/config
```

或者使用`format`参数来指定返回的配置格式

```text
https://127.0.0.1/manage/config?format=yaml
```

#### 报告接口

使用`report`命令行参数来启用报告接口的路由并且开始收集报告数据，需要先启用`manage`命令行参数

报告接口详情参见[Report API](/server/README.md#report-api)

```text
https://127.0.0.1/manage/report
```

### 日志

默认情况下会启用文件日志与控制台日志，你可以将`log_file`命令行参数设置为`false`来禁用文件日志

使用`log_level`命令行参数设置日志的等级，默认级别是`INFO`，可选项为：`DEBUG=0` `INFO=1` `WARN=2` `ERROR=3`

使用`log_dir`命令行参数来设置日志文件目录，默认为`./logs/`

使用`log_flush`命令行参数来设置自动刷新日志到文件中，默认启用

使用`log_flush_interval`命令行参数设置自动刷新日志到文件中的频率，默认为`3s`

使用`log_event`命令行参数启用事件日志，所有事件都会记录到文件中，默认为禁用

使用`log_sample_rate`命令行参数设置采样日志的采样率，取值范围为0到1，默认值为`1`

```bash
# 在"本地磁盘"模式下设置日志信息
$ gofs -source=./source -dest=./dest -log_file -log_level=0 -log_dir="./logs/" -log_flush -log_flush_interval=3s -log_event
```

### 使用配置文件

如果需要的话，你可以使用配置文件来代替所有的命令行参数，当前支持`json`和`yaml`格式

所有的配置字段名称跟命令行参数一样，你可以参考[配置示例](/conf/example)或者[配置接口](#配置接口)的响应结果

```bash
$ gofs -conf=./gofs.yaml
```

### 校验和

你可以使用`checksum`命令行参数来计算并打印文件的校验和

`chunk_size`和`checkpoint_count`命令行参数在这里同在[本地磁盘](#本地磁盘)中一样有效

```bash
$ gofs -source=./gofs -checksum
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

### 关于信息

```bash
$ gofs -about
```