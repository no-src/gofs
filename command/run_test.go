package command

import (
	"testing"
)

func TestRun(t *testing.T) {
	testCases := []commandCase{
		{"echo", run{Run: "echo hello world > run_echo.txt"}},
		{"mkdir", run{Run: "mkdir run_workspace"}},
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
