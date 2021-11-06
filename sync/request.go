package sync

type Request struct {
	Path    string
	Action  Action
	BaseUrl string
	// IsDir
	// 1 is dir
	// 0 not dir
	// -1 unknown
	IsDir int
	Size  int64
	Hash  string
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
