package version

import (
	"github.com/no-src/log"
)

const VERSION = "v0.2.2"

func PrintVersionInfo() {
	log.Log("gofs version %s", VERSION)
}
