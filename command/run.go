package command

type run struct {
	Run string `yaml:"run"`
}

func (c run) Exec() error {
	return nil
}
