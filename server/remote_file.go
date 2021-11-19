package server

type RemoteFile struct {
	Path  string `json:"path"`
	IsDir bool   `json:"is_dir"`
	Size  int64  `json:"size"`
	// crate time, unix sec
	CTime int64 `json:"c_time"`
	// last access time, unix sec
	ATime int64 `json:"a_time"`
	// last modify time, unix sec
	MTime int64 `json:"m_time"`
}
