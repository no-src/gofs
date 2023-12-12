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
	// MinIO the MinIO data source
	MinIO
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
	case MinIO:
		return "MinIO"
	default:
		return "Unknown"
	}
}
