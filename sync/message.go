package sync

import "github.com/no-src/gofs/contract"

// Message a message of the remote file change
type Message struct {
	contract.Status
	// Path the file path
	Path string `json:"path"`
	// Action the action of file change
	Action Action `json:"action"`
	// BaseUrl the base url of file server
	BaseUrl string `json:"base_url"`
	// IsDir is a dir the path
	IsDir contract.FsDirValue `json:"is_dir"`
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

type Action int

const (
	UnknownAction Action = iota
	// CreateAction the action of create a file
	CreateAction
	// WriteAction the action of write data to the file
	WriteAction
	// RemoveAction the action of remove the file
	RemoveAction
	// RenameAction the action of rename the file
	RenameAction
	// ChmodAction the action of change the file mode
	ChmodAction
)
