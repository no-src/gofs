package command

type mv struct {
	Source string `yaml:"source"`
	Dest   string `yaml:"dest"`
}

func (c mv) Exec() error {
	return nil
}
