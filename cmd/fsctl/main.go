package main

import (
	"flag"
	"os"

	"github.com/no-src/gofs/command"
	"github.com/no-src/gofs/internal/version"
	"github.com/no-src/log"
)

func main() {
	if c := run(); c != 0 {
		os.Exit(c)
	}
}

const errCode = 1

func run() (code int) {
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
		code = errCode
		return
	}

	if err := command.Exec(conf); err != nil {
		code = errCode
		log.Error(err, "execute commands failed")
	} else {
		log.Info("execute commands successfully")
	}
	return
}
