package sync

type Request struct {
	Path    string `json:"path"`
	Action  Action `json:"action"`
	BaseUrl string `json:"base_url"`
	// IsDir
	// 1 is dir
	// 0 not dir
	// -1 unknown
	IsDir int    `json:"is_dir"`
	Size  int64  `json:"size"`
	Hash  string `json:"hash"`
	// crate time, unix sec
	CTime int64 `json:"c_time"`
	// last access time, unix sec
	ATime int64 `json:"a_time"`
	// last modify time, unix sec
	MTime int64 `json:"m_time"`
}

type Action int

const (
	UnknownAction Action = iota
	CreateAction
	WriteAction
	RemoveAction
	RenameAction
	ChmodAction
)
