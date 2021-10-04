package core

// VFS virtual file system
type VFS struct {
	path   string
	fsType VFSType
}

// Path file path
func (vfs *VFS) Path() string {
	return vfs.path
}

// Type file system type
func (vfs *VFS) Type() VFSType {
	return vfs.fsType
}

// IsDisk is local file system
func (vfs *VFS) IsDisk() bool {
	return vfs.fsType == Disk
}

func NewDiskVFS(path string) VFS {
	vfs := VFS{
		fsType: Disk,
		path:   path,
	}
	return vfs
}

func NewEmptyVFS() VFS {
	vfs := VFS{
		fsType: Unknown,
	}
	return vfs
}
