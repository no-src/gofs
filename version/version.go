package version

import (
	"fmt"
	"strings"

	"github.com/no-src/gofs"
	"github.com/no-src/log"
)

// VERSION the current program version info
const VERSION = "v0.5.2"

// PrintVersion print the current version info, and append the commit info if the commit file is not empty
func PrintVersion(name string) {
	v := fmt.Sprintf("%s version %s", name, VERSION)
	if commit := strings.TrimSpace(gofs.Commit); len(commit) > 0 {
		v += fmt.Sprintf("\ngit commit %s", commit)
	}
	if goVersion := strings.TrimSpace(gofs.GoVersion); len(goVersion) > 0 {
		v += fmt.Sprintf("\ngo version %s", goVersion)
	}
	log.Log(v)
}
