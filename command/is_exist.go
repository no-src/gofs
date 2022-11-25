package command

type isExist struct {
	Source string `yaml:"source"`
	Expect bool   `yaml:"expect"`
}

func (c isExist) Exec() error {
	return nil
}
