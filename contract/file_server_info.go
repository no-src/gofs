package contract

type FileServerInfo struct {
	Status
	ServerAddr string `json:"server_addr"`
	SrcPath    string `json:"src_path"`
	TargetPath string `json:"target_path"`
	QueryAddr  string `json:"query_addr"`
}
