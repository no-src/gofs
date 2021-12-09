package sync

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
