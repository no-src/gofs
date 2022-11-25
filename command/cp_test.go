package command

import "testing"

func TestCp_ReturnError(t *testing.T) {
	testCases := []commandCase{
		{"open source error", cp{Source: "./not_exist.go", Dest: "./cp/not_exist.go.bak"}},
		{"create dest file error", cp{Source: "./cp_test.go", Dest: "./cp/not_exist.go.bak"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cmd.Exec()
			if err == nil {
				t.Errorf(testExecReturnErrorFailedMessage)
			}
		})
	}
}
