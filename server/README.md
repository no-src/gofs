# Server

## API List

| Name                                  | Route          | Method    | Remark       |
|---------------------------------------|----------------|-----------|--------------|
| Navigation Page                       | /              | GET       |              |
| Login Page                            | /login/index   | GET       |              |
| User Sign In API                      | /signin        | POST      |              |
| Source File Server                    | /source/       | GET       |              |
| DestPath File Server                  | /dest/         | GET       |              |
| [File Query API](#file-query-api)     | /query         | GET       |              |
| [File Push API](#file-push-api)       | /w/push        | POST      |              |
| PProf API                             | /manage/pprof  | GET       |              |
| Config API                            | /manage/config | GET       |              |

### File Query API

Support query source or dest path from [File Server](/README.md#file-server).

#### Request

##### Method

`GET`

##### Parameter

Request field description:

- `path` query file path, for example `path=source`
- `need_hash` return file hash or not, `1` or `0`, default is `0`

##### Example

Go to query the source path and return all file hash values.

```text
http://127.0.0.1/query?path=source&need_hash=1
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
      "path": "hello_gofs.txt",
      "is_dir": 0,
      "size": 11,
      "hash": "5eb63bbbe01eeed093cb22bb8f5acdc3",
      "c_time": 1642731076,
      "a_time": 1642731088,
      "m_time": 1642731088
    },
    {
      "path": "resource",
      "is_dir": 1,
      "size": 0,
      "hash": "",
      "c_time": 1642731096,
      "a_time": 1642731102,
      "m_time": 1642731102
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