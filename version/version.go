package version

import (
	"fmt"
	"github.com/no-src/gofs"
	"github.com/no-src/log"
	"strings"
)

// VERSION the current program version info
const VERSION = "v0.2.4"

// PrintVersion print the current version info, and append the commit info if the commit file is not empty
func PrintVersion() {
	v := fmt.Sprintf("gofs version %s", VERSION)
	if commit := strings.TrimSpace(gofs.Commit); len(commit) > 0 {
		v += fmt.Sprintf("\ngit commit %s", commit)
	}
	log.Log(v)
}
