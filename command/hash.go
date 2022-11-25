package command

type hash struct {
	Algorithm string `yaml:"algorithm"`
	Source    string `yaml:"source"`
	Expect    string `yaml:"expect"`
}

func (c hash) Exec() error {
	return nil
}
