package sync

type emptySync struct {
	baseSync
}

// NewEmptySync create a emptySync instance
func NewEmptySync(opt Option) (s Sync, err error) {
	// the fields of option
	source := opt.Source
	dest := opt.Dest
	logger := opt.Logger

	if source.IsEmpty() {
		return nil, errSourceNotFound
	}
	if dest.IsEmpty() {
		return nil, errDestNotFound
	}

	s = &emptySync{
		baseSync: newBaseSync(source, dest, logger),
	}
	return s, nil
}

func (s *emptySync) Create(path string) error {
	return nil
}

func (s *emptySync) Symlink(oldname, newname string) error {
	return nil
}

func (s *emptySync) Write(path string) error {
	return nil
}

func (s *emptySync) Remove(path string) error {
	return nil
}

func (s *emptySync) Rename(path string) error {
	return nil
}

func (s *emptySync) Chmod(path string) error {
	return nil
}

func (s *emptySync) IsDir(path string) (bool, error) {
	return false, nil
}

func (s *emptySync) SyncOnce(path string) error {
	return nil
}
