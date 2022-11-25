package command

type isDir struct {
	Source string `yaml:"source"`
	Expect bool   `yaml:"expect"`
}

func (c isDir) Exec() error {
	return nil
}
