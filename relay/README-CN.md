# 中继模式

[English](README.md) | 简体中文

如果你需要在两个无法直接相连的设备之间同步文件，可以使用反向代理作为中继服务器来实现，如[frp](https://github.com/fatedier/frp)、[ngrok](https://ngrok.com)等等

## 端口

### HTTP(S)

如果你显示或隐式的使用了[Web文件服务器](/README-CN.md#web文件服务器)，则需要开放指定的HTTP(S)端口，默认为`443`，实际值参见`server_addr`命令行参数的设定

### TCP

如果你启用了[远程磁盘服务端](/README-CN.md#远程磁盘服务端)或者[远程推送服务端](/README-CN.md#远程推送服务端)，则需要开放指定的TCP端口，默认为`8105`，实际值参见`source`
命令行参数中指定的[端口号](/README-CN.md#端口号)

## 反向代理

使用反向代理，将上述提到的端口服务暴露出去

### frp

使用[frp](https://github.com/fatedier/frp)实现在两个无法直接相连的设备之间进行文件同步

#### 服务端

服务端配置文件`frps.ini`内容如下

```text
[common]
bind_port = 7000
# HTTP REVERSE PROXY
vhost_http_port = 7001
# HTTPS REVERSE PROXY
vhost_https_port = 7002
```

启动服务端

```bash
$ ./frps -c ./frps.ini
```

#### 客户端

客户端配置文件`frpc.ini`内容如下

```text
[common]
server_addr = {YOUR FRP SERVER ADDRESS}
server_port = 7000

[gofs-https]
type = https
local_port = 443
custom_domains = {YOUR CUSTOM DOMAINS}

[gofs-tcp]
type = tcp
local_ip = 127.0.0.1
local_port = 8105
# TCP REVERSE PROXY
remote_port = 7003
```

启动客户端

```bash
$ ./frpc -c ./frpc.ini
```

### ngrok

使用[ngrok](https://ngrok.com)实现在两个无法直接相连的设备之间进行文件同步

#### 客户端

```bash
# HTTP(S) REVERSE PROXY
$ ngrok http https://127.0.0.1 --authtoken={YOUR NGROK TOEKN}

# TCP REVERSE PROXY
$ ngrok tcp 8105 --authtoken={YOUR NGROK TOEKN}
```

## 使用反向代理进行中继

以[远程磁盘服务端](/README-CN.md#远程磁盘服务端)与[远程磁盘客户端](/README-CN.md#远程磁盘客户端)为例，你只需要修改`source`命令行参数中对应的地址为反向代理最终暴露出来的地址即可

```bash
# 远程磁盘服务端
$ gofs -source="rs://127.0.0.1:8105?mode=server&local_sync_disabled=true&path=./source&fs_server={YOUR HTTP(S) REVERSE PROXY ADDRESS}" -dest=./dest -users="gofs|password|r" -tls_cert_file=cert.pem -tls_key_file=key.pem -token_secret=mysecret_16bytes

# 远程磁盘客户端
$ gofs -source="rs://{YOUR TCP REVERSE PROXY ADDRESS}" -dest=./dest -users="gofs|password" -tls_cert_file=cert.pem
```