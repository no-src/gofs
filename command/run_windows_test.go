package command

import (
	"testing"
)

func TestRun_Shell(t *testing.T) {
	testCases := []commandCase{
		{"cmd echo", run{Shell: "cmd", Run: "echo hello from cmd > run_echo_cmd.txt"}},
		{"cmd mkdir", run{Shell: "cmd", Run: "mkdir run_workspace_cmd"}},
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
