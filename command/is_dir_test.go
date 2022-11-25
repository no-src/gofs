package command

import "testing"

func TestIsDir_ReturnError(t *testing.T) {
	testCases := []commandCase{
		{"path is not exist", isDir{Source: "./not_exist.go"}},
		{"not match the expectation", isDir{Source: "./is_dir_test.go", Expect: true}},
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
