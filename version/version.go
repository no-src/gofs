package version

import "fmt"

const VERSION = "v0.0.1"

func PrintVersionInfo() {
	fmt.Printf("gofs %s\n", VERSION)
}
