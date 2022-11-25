package command

type isEmpty struct {
	Source string `yaml:"source"`
	Expect bool   `yaml:"expect"`
}

func (c isEmpty) Exec() error {
	return nil
}
