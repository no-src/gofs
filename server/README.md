# Server

## API List

| Name                                  | Route          | Method    | Remark       |
| ------------------------------------- | ---------------| ----------| -------------|
| Navigation Page                       | /              |    GET    |              |
| Login Page                            | /login/index   |    GET    |              |
| User Sign In API                      | /signin        |    POST   |              |
| Src File Server                       | /src/          |    GET    |              |
| Target File Server                    | /target/       |    GET    |              |
| [File Query API](#file-query-api)     | /query         |    GET    |              |
| PProf API                             | /debug/pprof   |    GET    |              |


### File Query API

Support query src or target path from [File Server](/README.md#file-server).

#### Method

`GET`

#### Parameter

- `path` query file path, for example `path=src`
- `need_hash` return file hash or not, `1` or `0`, default is `0`

#### Response

For example:

```text
http://127.0.0.1/query?path=src&need_hash=1
```

Response field description:

- `code` status code,`0` means success
- `message` response status description
- `data` response data
- `path` file path
- `is_dir` is directory or not, `1` or `0`
- `size` file size of bytes, directory is always `0`
- `hash` return file hash value if set `need_hash=1`
- `c_time` file create time
- `a_time` file last access time
- `m_time` file last modify time

An example response for query api.

```json
{
	"code": 0,
	"message": "success",
	"data": [{
		"path": "hello-gofs.txt",
		"is_dir": 0,
		"size": 11,
		"hash": "5eb63bbbe01eeed093cb22bb8f5acdc3",
		"c_time": 1642731076,
		"a_time": 1642731088,
		"m_time": 1642731088
	}, {
		"path": "resource",
		"is_dir": 1,
		"size": 0,
		"hash": "",
		"c_time": 1642731096,
		"a_time": 1642731102,
		"m_time": 1642731102
	}]
}
```