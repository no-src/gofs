# Relay

English | [简体中文](README-CN.md)

If you need to synchronize files between two devices that are unable to establish a direct connection, you can use a
reverse proxy as a relay server. For example: [frp](https://github.com/fatedier/frp), [ngrok](https://ngrok.com) etc.

## Ports

### HTTP(S)

If you explicitly or implicitly use the [File Server](/README.md#file-server), you will need to open the specified
HTTP(S) port, which defaults to `443`. See the `server_addr` flag for the actual value.

### TCP

If you enable the [Remote Disk Server](/README.md#remote-disk-server)
or [Remote Push Server](/README.md#remote-push-server)，you will need to open the specified TCP port, which defaults
to`8105`. See the [Port](/README.md#port) specified in the `source` flag for the actual value.

## Reverse Proxy

Use the reverse proxy to expose the port services mentioned above.

### frp

Use [frp](https://github.com/fatedier/frp) to achieve file synchronization between two devices that are unable to
establish a direct connection.

#### Server

The server config file `frps.ini` is as follows.

```text
[common]
bind_port = 7000
# HTTP REVERSE PROXY
vhost_http_port = 7001
# HTTPS REVERSE PROXY
vhost_https_port = 7002
```

Starting the server.

```bash
$ ./frps -c ./frps.ini
```

#### Client

The client config file `frpc.ini` is as follows.

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

Starting the client.

```bash
$ ./frpc -c ./frpc.ini
```

### ngrok

Use [ngrok](https://ngrok.com) to achieve file synchronization between two devices that are unable to establish a direct
connection.

#### Client

```bash
# HTTP(S) REVERSE PROXY
$ ngrok http https://127.0.0.1 --authtoken={YOUR NGROK TOEKN}

# TCP REVERSE PROXY
$ ngrok tcp 8105 --authtoken={YOUR NGROK TOEKN}
```

## Relaying With Reverse Proxy

Take the [Remote Disk Server](/README.md#remote-disk-server) and [Remote Disk Client](/README.md#remote-disk-client) as
an example, you simply change the address in the `source` flag to the address that the reverse proxy will eventually
expose.

```bash
# remote disk server
$ gofs -source="rs://127.0.0.1:8105?mode=server&local_sync_disabled=true&path=./source&fs_server={YOUR HTTP(S) REVERSE PROXY ADDRESS}" -dest=./dest -users="gofs|password|r" -tls_cert_file=cert.pem -tls_key_file=key.pem -token_secret=mysecret_16bytes

# remote disk client
$ gofs -source="rs://{YOUR TCP REVERSE PROXY ADDRESS}" -dest=./dest -users="gofs|password" -tls_cert_file=cert.pem
```