package main

import (
	"flag"

	"github.com/no-src/gofs/command"
	"github.com/no-src/gofs/version"
	"github.com/no-src/log"
)

func main() {
	run()
}

func run() {
	defer log.Close()

	var printVersion bool
	var conf string
	flag.BoolVar(&printVersion, "v", false, "print the version info")
	flag.StringVar(&conf, "conf", "", "the path of config file")
	flag.Parse()

	if printVersion {
		version.PrintVersion("fsctl")
		return
	}

	if len(conf) == 0 {
		log.Info("please specify the config file by -conf flag")
		return
	}

	if err := command.Exec(conf); err != nil {
		log.Error(err, "execute commands failed")
	} else {
		log.Info("execute commands successfully")
	}
}
