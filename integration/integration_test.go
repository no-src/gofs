//go:build integration_test

package integration

import (
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
	return cmd.RunWithConfigFile(path)
}
