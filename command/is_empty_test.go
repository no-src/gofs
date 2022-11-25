package command

import "testing"

func TestIsEmpty_ReturnError(t *testing.T) {
	testCases := []commandCase{
		{"path is not exist", isEmpty{Source: "./not_exist.go"}},
		{"not match the expectation", isEmpty{Source: "./is_empty_test.go", Expect: true}},
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
