package command

type run struct {
	Run   string `yaml:"run"`
	Shell string `yaml:"shell"`
}

func (c run) Name() string {
	return "run"
}

func init() {
	registerCommand("run", func(a Action) (Command, error) {
		return parse[run](a)
	})
}
