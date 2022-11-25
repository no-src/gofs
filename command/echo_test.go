package command

import "testing"

func TestEcho_ReturnError(t *testing.T) {
	testCases := []commandCase{
		{"open source error", echo{Source: "./echo/not_exist.go.bak"}},
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
