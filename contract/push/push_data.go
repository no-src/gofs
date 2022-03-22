package push

import (
	"github.com/no-src/gofs/action"
	"github.com/no-src/gofs/contract"
)

// PushData the request data of the push api
type PushData struct {

	// Action the action of file change
	Action action.Action `json:"action"`

	// PushAction the push action of comparing or writing to the file
	PushAction PushAction `json:"push_action"`

	// FileInfo the basic file info
	FileInfo contract.FileInfo `json:"file_info"`

	// Chunk the basic file chunk info
	Chunk contract.Chunk `json:"chunk"`
}
