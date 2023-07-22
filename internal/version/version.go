package version

import (
	_ "embed"
	"fmt"
	"runtime"
	"strings"

	"github.com/no-src/log"
)

// VERSION the current program version info
const VERSION = "v0.7.2"

// Commit the commit file records the last commit hash value, used by release
//
//go:embed commit
var Commit string

// PrintVersion print the current version info, and append the commit info if the commit file is not empty
func PrintVersion(name string) {
	v := fmt.Sprintf("%s version %s", name, VERSION)
	if commit := strings.TrimSpace(Commit); len(commit) > 0 {
		v += fmt.Sprintf("\ngit commit %s", commit)
	}
	v += fmt.Sprintf("\ngo version %s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	log.Log(v)
}
