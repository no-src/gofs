package command

import (
	"os/exec"
)

func (c run) Exec() error {
	shell := "cmd"
	if len(c.Shell) > 0 {
		shell = c.Shell
	}
	return exec.Command(shell, "/c", c.Run).Run()
}
