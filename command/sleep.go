package command

import "time"

type sleep struct {
	Sleep time.Duration `yaml:"sleep"`
}

func (c sleep) Exec() error {
	time.Sleep(c.Sleep)
	return nil
}

func (c sleep) Name() string {
	return "sleep"
}

func init() {
	registerCommand("sleep", func(a Action) (Command, error) {
		return parse[sleep](a)
	})
}
