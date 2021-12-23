package contract

type FileInfo struct {
	// Path the file path
	Path string `json:"path"`
	// IsDir is a dir the path
	IsDir FsDirValue `json:"is_dir"`
	// Size the size of path for bytes
	Size int64 `json:"size"`
	// Hash calculate the path hash value, if the path is a file
	Hash string `json:"hash"`
	// CTime create time, unix sec
	CTime int64 `json:"c_time"`
	// ATime last access time, unix sec
	ATime int64 `json:"a_time"`
	// MTime last modify time, unix sec
	MTime int64 `json:"m_time"`
}
