package command

import "time"

type sleep struct {
	Sleep time.Duration `yaml:"sleep"`
}

func (c sleep) Exec() error {
	time.Sleep(c.Sleep)
	return nil
}
