package version

import (
	"github.com/no-src/log"
)

const VERSION = "v0.0.3"

func PrintVersionInfo() {
	log.Log("gofs version %s", VERSION)
}
