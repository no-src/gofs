package version

import (
	"github.com/no-src/log"
)

const VERSION = "v0.2.1"

func PrintVersionInfo() {
	log.Log("gofs version %s", VERSION)
}
