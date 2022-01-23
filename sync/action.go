package sync

type Action int

const (
	// UnknownAction the unknown file operation
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

// String return the action description name
func (action Action) String() string {
	desc := ""
	switch action {
	case CreateAction:
		desc = "Create"
		break
	case WriteAction:
		desc = "Write"
		break
	case RemoveAction:
		desc = "Remove"
		break
	case RenameAction:
		desc = "Rename"
		break
	case ChmodAction:
		desc = "Chmod"
		break
	default:
		desc = "Unknown"
		break
	}
	return desc
}
