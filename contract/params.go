package contract

const (
	// FsDir is dir, see FsDirValue
	FsDir = "dir"
	// FsSize file size, bytes
	FsSize = "size"
	// FsHash file hash value
	FsHash = "hash"
	// FsHashValues the hash value of the entire file and first chunk and some checkpoints
	FsHashValues = "hash_values"
	// FsCtime file creation time
	FsCtime = "ctime"
	// FsAtime file last access time
	FsAtime = "atime"
	// FsMtime file last modify time
	FsMtime = "mtime"
	// FsPath file path
	FsPath = "path"
	// FsNeedHash return file hash or not
	FsNeedHash = "need_hash"
	// FsNeedCheckpoint return file checkpoint hash or not
	FsNeedCheckpoint = "need_checkpoint"
)

const (
	// ParamValueTrue the parameter value means true
	ParamValueTrue = "1"
	// ParamValueFalse the parameter value means false
	ParamValueFalse = "0"
	// FsNeedHashValueTrue the optional value of the FsNeedHash parameter, that means let file server return file hash value
	FsNeedHashValueTrue = ParamValueTrue
)
