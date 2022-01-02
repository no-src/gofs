package version

import (
	"fmt"
	"github.com/no-src/gofs"
	"github.com/no-src/log"
)

const VERSION = "v0.2.2"

func PrintVersionInfo() {
	v := fmt.Sprintf("gofs version %s", VERSION)
	commit, err := gofs.Version.ReadFile("version/commit")
	if err == nil && len(commit) > 0 {
		v += fmt.Sprintf("\ngit commit %s", string(commit))
	}
	log.Log(v)
}
