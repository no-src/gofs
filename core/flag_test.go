package core

import (
	"flag"
	"os"
)

var (
	testCommandLine = NewFlagSet(os.Args[0], flag.ExitOnError)
)
