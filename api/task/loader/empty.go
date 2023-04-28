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
