package monitor

type sftpPullClientMonitor struct {
	driverPullClientMonitor
}

// NewSftpPullClientMonitor create an instance of sftpPullClientMonitor to pull the files from sftp server
func NewSftpPullClientMonitor(opt Option) (m Monitor, err error) {
	m = &sftpPullClientMonitor{
		driverPullClientMonitor: driverPullClientMonitor{
			baseMonitor: newBaseMonitor(opt),
		},
	}
	return m, nil
}
