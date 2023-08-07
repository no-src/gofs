package contract

import "github.com/no-src/gofs/util/hashutil"

// FileInfo the basic file info description
type FileInfo struct {
	// Path the file path
	Path string `json:"path"`
	// IsDir is a dir the path
	IsDir FsDirValue `json:"is_dir"`
	// Size the size of path for bytes
	Size int64 `json:"size"`
	// Hash calculate the path hash value, if the path is a file
	Hash string `json:"hash"`
	// HashValues the hash value of the entire file and first chunk and some checkpoints
	HashValues hashutil.HashValues `json:"hash_values"`
	// CTime creation time, unix sec
	CTime int64 `json:"c_time"`
	// ATime last access time, unix sec
	ATime int64 `json:"a_time"`
	// MTime last modify time, unix sec
	MTime int64 `json:"m_time"`
	// LinkTo link to the real file
	LinkTo string `json:"link_to"`
}
