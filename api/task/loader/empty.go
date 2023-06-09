package loader

type emptyLoader struct {
}

func newEmptyLoader() Loader {
	return &emptyLoader{}
}

func (loader *emptyLoader) LoadConfig() (c *TaskConfig, err error) {
	return &TaskConfig{}, nil
}

func (loader *emptyLoader) LoadContent(conf string) (content string, err error) {
	return "", nil
}

func (loader *emptyLoader) SaveConfig(c *TaskConfig) error {
	return nil
}

func (loader *emptyLoader) SaveContent(conf string, content string) error {
	return nil
}

func (loader *emptyLoader) Close() error {
	return nil
}
