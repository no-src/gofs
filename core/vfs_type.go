package core

// VFSType the file data source type
type VFSType int

const (
	// Unknown the unknown file data source
	Unknown VFSType = iota
	// Disk the local disk file system data source
	Disk
	// RemoteDisk the remote disk file system data source
	RemoteDisk
	// FTP the ftp data source
	FTP
	// SFTP the sftp data source
	SFTP
	// NetworkDisk the network disk data source
	NetworkDisk
	// SharedFolder the shared folder data source
	SharedFolder
)

// String return the VFSType name
func (t VFSType) String() string {
	switch t {
	case Disk:
		return "Disk"
	case RemoteDisk:
		return "RemoteDisk"
	case FTP:
		return "FTP"
	case SFTP:
		return "SFTP"
	case NetworkDisk:
		return "NetworkDisk"
	case SharedFolder:
		return "SharedFolder"
	default:
		return "Unknown"
	}
}
