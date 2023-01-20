//go:build linux || darwin

package command

import "testing"

func TestRun_Shell(t *testing.T) {
	testCases := []commandCase{
		{"bash echo", run{Shell: "bash", Run: "echo hello from bash > run_echo_bash.txt"}},
		{"sh echo", run{Shell: "sh", Run: "echo hello from sh > run_echo_sh.txt"}},
		{"bash mkdir", run{Shell: "bash", Run: "mkdir run_workspace_bash"}},
		{"sh mkdir", run{Shell: "sh", Run: "mkdir run_workspace_sh"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cmd.Exec()
			if err != nil {
				t.Errorf("run command error, err=%v", err)
			}
		})
	}
}
