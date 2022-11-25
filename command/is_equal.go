package command

type isEqual struct {
	Source       string `yaml:"source"`
	Dest         string `yaml:"dest"`
	Expect       bool   `yaml:"expect"`
	Size         int64  `yaml:"size"`
	MustNonEmpty bool   `yaml:"must-non-empty"`
}

func (c isEqual) Exec() error {
	return nil
}
