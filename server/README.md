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

- `code` status code,`1` means success
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

- `file_info` basic push file info
    - `action` the action of file change, Create(1) Write(2) Remove(3) Rename(4) Chmod(5)
    - `path` file path
    - `is_dir` is directory or not, `1` or `0`
    - `size` file size of bytes, directory is always `0`
    - `hash` return file hash value if set `need_hash=1`
    - `c_time` file create time
    - `a_time` file last access time
    - `m_time` file last modify time
- `offset` the offset relative to the origin of the file, `-1` means to compare file size and hash value only, `0` means
  the first chunk or only one chunk
- `up_file` the field name of upload file

##### Example

Upload a file to the remote push server.

```text
POST https://127.0.0.1/w/push HTTP/1.1
Host: 127.0.0.1
User-Agent: Go-http-client/1.1
Content-Length: 646
Content-Type: multipart/form-data; boundary=af3294e968a2357d7cd21f809d3508ef96e1db1621bdc7bd1b321160676b
Cookie: session_id=MTY0NjM3MDU2OXxOd3dBTkVGUldrUlhSa3RCUmsxVFVFcEtOMHhHUjFKRVExRmFWMWhTVjFkSlNVWk1XVGRZU0RWS05VVlFRbGswVDB0U1UwRkJTRUU9fFrO-f0mlkXZFQvCJGv_ufJTqgmmrEQPoTLKFhcWG5_D
Accept-Encoding: gzip

--af3294e968a2357d7cd21f809d3508ef96e1db1621bdc7bd1b321160676b
Content-Disposition: form-data; name="file_info"

{"path":"hello_gofs.txt","is_dir":0,"size":5,"hash":"5d41402abc4b2a76b9719d911017c592","c_time":1646370569,"a_time":1646370572,"m_time":1646287764,"action":2}
--af3294e968a2357d7cd21f809d3508ef96e1db1621bdc7bd1b321160676b
Content-Disposition: form-data; name="offset"

0
--af3294e968a2357d7cd21f809d3508ef96e1db1621bdc7bd1b321160676b
Content-Disposition: form-data; name="up_file"; filename="hello_gofs.txt"
Content-Type: application/octet-stream

hello
--af3294e968a2357d7cd21f809d3508ef96e1db1621bdc7bd1b321160676b--
```

#### Response

##### Parameter

Response field description:

- `code` status code,`1` means success
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