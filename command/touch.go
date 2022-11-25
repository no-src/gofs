package command

type touch struct {
	Source string `yaml:"source"`
}

func (c touch) Exec() error {
	return nil
}
