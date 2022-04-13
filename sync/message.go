package sync

import (
	"github.com/no-src/gofs/action"
	"github.com/no-src/gofs/contract"
)

// Message a message of the remote file change
type Message struct {
	contract.Status
	contract.FileInfo

	// Action the action of file change
	Action action.Action `json:"action"`
	// BaseUrl the base url of file server
	BaseUrl string `json:"base_url"`
}
