package action

import "strconv"

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
	// maxAction the max boundary value of Action, it is an invalid value
	maxAction
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
	case UnknownAction:
		desc = "Unknown"
		break
	default:
		desc = "Invalid"
		break
	}
	return desc
}

func (action Action) Int() int {
	return int(action)
}

func (action Action) Valid() Action {
	if action >= maxAction || action <= UnknownAction {
		return UnknownAction
	}
	return action
}

func ParseActionFromString(action string) Action {
	i, err := strconv.Atoi(action)
	if err != nil {
		return UnknownAction
	}
	return ParseAction(i)
}

func ParseAction(action int) Action {
	a := Action(action)
	return a.Valid()
}
