package contract

// FileServerInfo the file server basic info
type FileServerInfo struct {
	Status
	// ServerAddr the server running address
	ServerAddr string `json:"server_addr"`
	// SrcPath the src base path of the file server
	SrcPath string `json:"src_path"`
	// TargetPath the target base path of the file server
	TargetPath string `json:"target_path"`
	// QueryAddr the query api path of the file server
	QueryAddr string `json:"query_addr"`
}
