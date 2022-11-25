package command

import "testing"

func TestIsExist_ReturnError(t *testing.T) {
	testCases := []commandCase{
		{"not match the expectation", isExist{Source: "./is_exist_test.go", Expect: false}},
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
