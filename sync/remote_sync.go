package sync

// NewRemoteSync auto create an instance of remoteServerSync or remoteClientSync according to source and dest
func NewRemoteSync(opt Option) (Sync, error) {
	if opt.Source.Server() {
		return NewRemoteServerSync(opt)
	}
	return NewRemoteClientSync(opt)
}
