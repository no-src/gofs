package command

type mkdir struct {
	Source string `yaml:"source"`
}

func (c mkdir) Exec() error {
	return nil
}
