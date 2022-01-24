package contract

// FileServerInfo the file server basic info
type FileServerInfo struct {
	Status
	// ServerAddr the server running address
	ServerAddr string `json:"server_addr"`
	// SrcPath the src base path of the file server
	SrcPath string `json:"src_path"`
	// DestPath the dest base path of the file server
	DestPath string `json:"dest_path"`
	// QueryAddr the query api path of the file server
	QueryAddr string `json:"query_addr"`
}
