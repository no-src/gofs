//go:build integration_test

package integration

import (
	"flag"
	"os"

	"github.com/no-src/gofs/cmd"
	"github.com/no-src/gofs/result"
)

func getRunConf(conf string) string {
	return "./testdata/conf/" + conf
}

func getTestConf(conf string) string {
	return "./testdata/test/" + conf
}

func runWithConfigFile(path string) result.Result {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	return cmd.RunWithConfigFile(path)
}
