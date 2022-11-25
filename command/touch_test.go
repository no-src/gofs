package command

import "testing"

func TestTouch_ReturnError(t *testing.T) {
	testCases := []commandCase{
		{"directory is not exist", touch{Source: "./touch/hello"}},
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
