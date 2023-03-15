package command

type run struct {
	Run   string `yaml:"run"`
	Shell string `yaml:"shell"`
}

func (c run) Name() string {
	return "run"
}

func init() {
	registerCommand[run]()
}
