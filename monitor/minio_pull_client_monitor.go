package monitor

type minIOPullClientMonitor struct {
	driverPullClientMonitor
}

// NewMinIOPullClientMonitor create an instance of minIOPullClientMonitor to pull the files from MinIO server
func NewMinIOPullClientMonitor(opt Option) (m Monitor, err error) {
	m = &minIOPullClientMonitor{
		driverPullClientMonitor: driverPullClientMonitor{
			baseMonitor: newBaseMonitor(opt),
		},
	}
	return m, nil
}
