package command

type rm struct {
	Source string `yaml:"source"`
}

func (c rm) Exec() error {
	return nil
}
