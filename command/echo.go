package command

type echo struct {
	Source string `yaml:"source"`
	Input  string `yaml:"input"`
	Append bool   `yaml:"append"`
}

func (c echo) Exec() error {
	return nil
}
