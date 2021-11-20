package version

import (
	"github.com/no-src/log"
)

const VERSION = "v0.1.0"

func PrintVersionInfo() {
	log.Log("gofs version %s", VERSION)
}
