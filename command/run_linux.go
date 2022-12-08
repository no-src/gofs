package command

import (
	"os/exec"
)

func (c run) Exec() error {
	shell := "sh"
	if len(c.Shell) > 0 {
		shell = c.Shell
	}
	return exec.Command(shell, "-c", c.Run).Run()
}
