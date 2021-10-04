package core

type VFSType int

const (
	Unknown VFSType = iota
	Disk
	FTP
	SFTP
	NetworkDisk
	SharedFolder
)

func (t VFSType) String() string {
	switch t {
	case Disk:
		return "Disk"
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
