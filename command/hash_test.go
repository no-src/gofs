package command

import "testing"

func TestHash_ReturnError(t *testing.T) {
	testCases := []commandCase{
		{"invalid algorithm", hash{Algorithm: "", Source: "./hash_test.go"}},
		{"path is not exist", hash{Algorithm: "md5", Source: "./not_exist.go"}},
		{"not match the expectation", hash{Algorithm: "md5", Source: "./hash_test.go", Expect: "123456"}},
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
