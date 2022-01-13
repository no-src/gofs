package version

import (
	"fmt"
	"github.com/no-src/gofs"
	"github.com/no-src/log"
	"strings"
)

const VERSION = "v0.2.3"

func PrintVersionInfo() {
	v := fmt.Sprintf("gofs version %s", VERSION)
	if commit := strings.TrimSpace(gofs.Commit); len(commit) > 0 {
		v += fmt.Sprintf("\ngit commit %s", commit)
	}
	log.Log(v)
}
