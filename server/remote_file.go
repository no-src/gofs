package server

type RemoteFile struct {
	Path  string
	IsDir bool
	Size  int64
}
