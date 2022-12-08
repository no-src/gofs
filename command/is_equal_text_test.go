package command

import (
	"os"
	"testing"
)

func TestIsEqualText_ReturnError(t *testing.T) {
	text, err := os.ReadFile("./is_equal_text_test.go")
	if err != nil {
		t.Errorf("read test file error, err=%v", err)
		return
	}
	testCases := []commandCase{
		{"source path is not exist", isEqualText{Source: "./not_exist.go", Dest: "./is_equal_text_test.go"}},
		{"size not equal", isEqualText{Source: "./is_equal_text_test.go", Dest: "hello", Expect: true}},
		{"not match the expectation", isEqualText{Source: "./is_equal_text_test.go", Dest: string(text), Expect: false}},
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
