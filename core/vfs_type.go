package core

type VFSType int

const (
	Unknown VFSType = iota
	Disk
	RemoteDisk
	FTP
	SFTP
	NetworkDisk
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
