package server

import "github.com/no-src/gofs/contract"

// Info file server info
type Info struct {
	contract.Status
	ServerAddr string `json:"server_addr"`
	SrcPath    string `json:"src_path"`
	TargetPath string `json:"target_path"`
	QueryAddr  string `json:"query_addr"`
}
