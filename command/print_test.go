package command

import "testing"

func TestPrint(t *testing.T) {
	testCases := []commandCase{
		{"print hello world", print{Input: "hello world"}},
		{"print empty line", print{Input: ""}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.cmd.Exec()
			if err != nil {
				t.Errorf("expect to get an nil error, but get %v", err)
			}
		})
	}
}
