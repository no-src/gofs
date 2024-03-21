package about

import (
	"fmt"
	"testing"
)

func TestPrintAbout(t *testing.T) {
	var s string
	PrintAbout(func(format string, args ...any) {
		s += fmt.Sprintf(format+"\n", args...)
	})
	fmt.Println(s)
}
