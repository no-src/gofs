package action

import "strconv"

// Action the action of file change
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
	case WriteAction:
		desc = "Write"
	case RemoveAction:
		desc = "Remove"
	case RenameAction:
		desc = "Rename"
	case ChmodAction:
		desc = "Chmod"
	case UnknownAction:
		desc = "Unknown"
	default:
		desc = "Invalid"
	}
	return desc
}

// Int return the int value of Action
func (action Action) Int() int {
	return int(action)
}

// Valid if the current Action is an invalid int value, return the UnknownAction
func (action Action) Valid() Action {
	if action >= maxAction || action <= UnknownAction {
		return UnknownAction
	}
	return action
}

// ParseActionFromString parse the string value to Action
func ParseActionFromString(action string) Action {
	i, err := strconv.Atoi(action)
	if err != nil {
		return UnknownAction
	}
	return ParseAction(i)
}

// ParseAction parse the int value to Action
func ParseAction(action int) Action {
	return Action(action).Valid()
}
