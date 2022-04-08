# Server

## API List

| Name                              | Route          | Method    | Remark       |
|-----------------------------------|----------------|-----------|--------------|
| Navigation Page                   | /              | GET       |              |
| Login Page                        | /login/index   | GET       |              |
| User Sign In API                  | /signin        | POST      |              |
| Source File Server                | /source/       | GET       |              |
| DestPath File Server              | /dest/         | GET       |              |
| [File Query API](#file-query-api) | /query         | GET       |              |
| [File Push API](#file-push-api)   | /w/push        | POST      |              |
| PProf API                         | /manage/pprof  | GET       |              |
| Config API                        | /manage/config | GET       |              |
| [Report API](#report-api)         | /manage/report | GET       |              |

### File Query API

Support query source or dest path from [File Server](/README.md#file-server).

#### Request

##### Method

`GET`

##### Parameter

Request field description:

- `path` query file path, for example `path=source`
- `need_hash` return file hash or not, `1` or `0`, default is `0`
- `need_checkpoint` return file checkpoint hash or not, `1` or `0`, default is `0`

##### Example

Query the source path and return all files and checkpoints hash values.

```text
http://127.0.0.1/query?path=source&need_hash=1&need_checkpoint=1
```

#### Response

##### Parameter

Response field description:

- `code` status code,`1` means success, all status codes see [Status Code](#status-code)
- `message` response status description
- `data` response data
    - `path` file path
    - `is_dir` is directory or not, `1` or `0`
    - `size` file size of bytes, directory is always `0`
    - `hash` return file hash value if set `need_hash=1`
    - `hash_values` return the hash value of the entire file and first chunk and some checkpoints if
      set `need_checkpoint=1`
        - `offset` the file data to calculate the hash value from zero to offset
        - `hash` the file checkpoint hash value
    - `c_time` file create time
    - `a_time` file last access time
    - `m_time` file last modify time

##### Example

Here is an example response:

```json
{
  "code": 1,
  "message": "success",
  "data": [
    {
      "path": "go1.18.linux-amd64.tar.gz",
      "is_dir": 0,
      "size": 141702072,
      "hash": "67622d8e307cb055b8ce2c2e03cb58cf",
      "hash_values": [
        {
          "offset": 1048576,
          "hash": "29d68359be77cdbe3d59791f0e367012"
        },
        {
          "offset": 13631488,
          "hash": "576b3746abb1f71bfbc794c750fb2c19"
        },
        {
          "offset": 27262976,
          "hash": "3a720d0a1a1c8abce3bf39cb7bc38507"
        },
        {
          "offset": 40894464,
          "hash": "946efd928cda9d0b20e9a74c4ba5db4f"
        },
        {
          "offset": 54525952,
          "hash": "eae47693bea9473c5ed859862685e4a7"
        },
        {
          "offset": 68157440,
          "hash": "ea700e864621fbee4a50d7b0e70b2e52"
        },
        {
          "offset": 81788928,
          "hash": "6c68022d11128ba256fb73a90b5232ef"
        },
        {
          "offset": 95420416,
          "hash": "ccfb4397caee4088b90e45d07829f1bf"
        },
        {
          "offset": 109051904,
          "hash": "d01e494ca88aa08d8a8a5279948db0b0"
        },
        {
          "offset": 122683392,
          "hash": "86bf50e8553674c26e2fbdb4477e44f4"
        },
        {
          "offset": 136314880,
          "hash": "51bba2e19d9519babccb0d5f4eedd7c3"
        },
        {
          "offset": 141702072,
          "hash": "67622d8e307cb055b8ce2c2e03cb58cf"
        }
      ],
      "c_time": 1649431872,
      "a_time": 1649431873,
      "m_time": 1647397031
    },
    {
      "path": "hello_gofs.txt",
      "is_dir": 0,
      "size": 11,
      "hash": "5eb63bbbe01eeed093cb22bb8f5acdc3",
      "hash_values": [
        {
          "offset": 11,
          "hash": "5eb63bbbe01eeed093cb22bb8f5acdc3"
        }
      ],
      "c_time": 1649431542,
      "a_time": 1649434237,
      "m_time": 1649434237
    },
    {
      "path": "resource",
      "is_dir": 1,
      "size": 0,
      "hash": "",
      "hash_values": null,
      "c_time": 1649431669,
      "a_time": 1649431898,
      "m_time": 1649431898
    }
  ]
}
```

### File Push API

Push the file changes to the [Remote Push Server](/README.md#remote-push-server).

#### Request

##### Method

`POST`

##### Parameter

Request field description:

- `push_data` the request data of push api, contains basic push file info and chunk info etc
    - `action` the action of file change, Create(1) Write(2) Remove(3) Rename(4) Chmod(5)
    - `push_action` the file upload action, CompareFile(1) CompareChunk(2) CompareFileAndChunk(3) Write(4) Truncate(5)
    - `file_info` basic push file info
        - `path` file path
        - `is_dir` is directory or not, `1` or `0`
        - `size` file size of bytes, directory is always `0`
        - `hash` file hash value
        - `c_time` file create time
        - `a_time` file last access time
        - `m_time` file last modify time
    - `chunk`
        - `offset` the offset relative to the origin of the file
        - `size` file chunk size of bytes, directory is always `0`
        - `hash` file chunk hash value
- `up_file` the field name of upload file or chunk

##### Example

Upload a file to the remote push server.

```text
POST https://127.0.0.1/w/push HTTP/1.1
Host: 127.0.0.1
User-Agent: Go-http-client/1.1
Content-Length: 632
Content-Type: multipart/form-data; boundary=d633162324641d10fc7ed0c03f2632807141c09b2b9e91a5b502b838fae7
Cookie: session_id=MTY0Nzk3NTMwMnxOd3dBTkU1UlMxb3lObE16VkZaVlZFMVdTa2d5VUV0TlNUZFdXbGRhVlVsSk4wNUhWVXRCVlV4RU0wdFlRVWRKUzFwQ1FqWlhTVUU9fD0l0l5GztC1TeOaR75R_dm90Z2c1q1X7xPPPMO2OPZl
Accept-Encoding: gzip

--d633162324641d10fc7ed0c03f2632807141c09b2b9e91a5b502b838fae7
Content-Disposition: form-data; name="push_data"

{"action":2,"push_action":4,"file_info":{"path":"hello_gofs.txt","is_dir":0,"size":5,"hash":"5d41402abc4b2a76b9719d911017c592","c_time":1647974350,"a_time":1647975313,"m_time":1647975313},"chunk":{"offset":0,"hash":"5d41402abc4b2a76b9719d911017c592","size":5}}
--d633162324641d10fc7ed0c03f2632807141c09b2b9e91a5b502b838fae7
Content-Disposition: form-data; name="up_file"; filename="hello_gofs.txt"
Content-Type: application/octet-stream

hello
--d633162324641d10fc7ed0c03f2632807141c09b2b9e91a5b502b838fae7--
```

#### Response

##### Parameter

Response field description:

- `code` status code,`1` means success, all status codes see [Status Code](#status-code)
- `message` response status description
- `data` response data

##### Example

Here is an example response:

```json
{
  "code": 1,
  "message": "success",
  "data": null
}
```

### Report API

Query the report data if you enable the `manage` and `report` flags.

#### Request

##### Method

`GET`

##### Example

```text
https://127.0.0.1/manage/report
```

#### Response

##### Parameter

Response field description:

- `code` status code,`1` means success, all status codes see [Status Code](#status-code)
- `message` response status description
- `data` response data
    - `pid` returns the process id of the caller
    - `ppid` returns the process id of the caller's parent
    - `go_os` is the running program's operating system target
    - `go_arch` is the running program's architecture target
    - `go_version` returns the Go tree's version string
    - `version` returns the version info of the gofs
    - `online` returns the client connection info that is online
        - `addr` the client connection address
        - `is_auth` whether the client is authorized
        - `username` the username of client
        - `perm` the permission of client
        - `connect_time` the connected time of client
        - `auth_time` the authorized time of client
        - `disconnect_time` the disconnected time of client
        - `life_time` the lifetime of a client, it is `0s` always that if the client is online
    - `offline` returns the client connection info that is offline, full fields see `online`
    - `events` returns some latest file change events
        - `name` the path of file change
        - `op` the operation of file change
        - `time` the time of file change
    - `event_stat` returns the statistical data of file change events
    - `api_stat` returns the statistical data of api access info
        - `access_count` all the api access count
        - `visitor_stat` the statistical data of visitors

##### Example

Here is an example response:

```json
{
  "code": 1,
  "message": "success",
  "data": {
    "pid": 94032,
    "ppid": 9268,
    "go_os": "windows",
    "go_arch": "amd64",
    "go_version": "go1.18",
    "version": "v0.4.0",
    "online": {
      "127.0.0.1:11993": {
        "addr": "127.0.0.1:11993",
        "is_auth": true,
        "username": "698d51a19d8a121c",
        "perm": "rwx",
        "connect_time": "2022-03-28 01:10:11",
        "auth_time": "2022-03-28 01:10:11",
        "disconnect_time": "1970-01-01 08:00:00",
        "life_time": "0s"
      }
    },
    "offline": [
      {
        "addr": "127.0.0.1:11887",
        "is_auth": true,
        "username": "698d51a19d8a121c",
        "perm": "rwx",
        "connect_time": "2022-03-28 01:08:46",
        "auth_time": "2022-03-28 01:08:46",
        "disconnect_time": "2022-03-28 01:10:06",
        "life_time": "1m20s"
      }
    ],
    "events": [
      {
        "name": "C:\\workspace\\hello_gofs.txt",
        "op": "WRITE",
        "time": "2022-03-28 01:10:01"
      }
    ],
    "event_stat": {
      "WRITE": 1
    },
    "api_stat": {
      "access_count": 14,
      "visitor_stat": {
        "127.0.0.1": 11,
        "192.168.0.106": 3
      }
    }
  }
}
```

## Status Code

All common response status code enums below.

- `0`   Unknown
- `1`   Success
- `-1`  Fail
- `-2`  Unauthorized
- `-3`  NotFound
- `-4`  NoPermission
- `-5`  ServerError
- `-6`  AccessDeny
- `-7`  NotModified
- `-8`  ChunkNotModified
- `-9`  Modified
- `-10` ChunkModified