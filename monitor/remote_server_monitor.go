package monitor

// NewRemoteServerMonitor create an instance of the fsNotifyMonitor
func NewRemoteServerMonitor(opt Option) (m Monitor, err error) {
	return NewFsNotifyMonitor(opt)
}
