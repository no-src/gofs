package main

import (
	"os"

	"github.com/no-src/gofs/cmd"
)

func main() {
	if err := cmd.Run().Wait(); err != nil {
		os.Exit(1)
	}
}
