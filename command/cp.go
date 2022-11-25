package command

type cp struct {
	Source string `yaml:"source"`
	Dest   string `yaml:"dest"`
}

func (c cp) Exec() error {
	return nil
}
