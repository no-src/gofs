package command

import "testing"

func TestIsEqual_ReturnError(t *testing.T) {
	testCases := []commandCase{
		{"invalid hash algorithm", isEqual{Source: "./is_equal_test.go", Dest: "./is_equal.go", Expect: false, Algorithm: "invalid"}},
		{"source path is not exist", isEqual{Source: "./not_exist.go", Dest: "./is_equal_test.go"}},
		{"dest path is not exist", isEqual{Source: "./is_equal_test.go", Dest: "./not_exist.go"}},
		{"zero source file", isEqual{Source: "../internal/version/commit", Dest: "./is_equal_test.go", Expect: true, MustNonEmpty: true}},
		{"zero dest file", isEqual{Source: "./is_equal_test.go", Dest: "../internal/version/commit", Expect: true, MustNonEmpty: true}},
		{"size not equal", isEqual{Source: "./is_equal_test.go", Dest: "./command_test.go", Expect: true}},
		{"not match the expectation", isEqual{Source: "./is_equal_test.go", Dest: "./is_equal_test.go", Expect: false}},
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
