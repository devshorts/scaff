package main

import (
	"os"

	"github.com/devshorts/scaff/scaff"
	"github.com/devshorts/scaff/scaff/file"
	"github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
)

func main() {
	opts := scaff.Opts{}

	parser := flags.NewParser(&opts, flags.Default)

	if _, e := parser.Parse(); e != nil {
		os.Exit(-1)
	}

	if file.Exists(opts.TargetDir) {
		logrus.Warnf("%s already exists, please remove it before re-running", opts.TargetDir)

		os.Exit(-1)
	}

	// load a config from the scaff config file in the template directory
	config := scaff.NewParser(opts.ScaffConfigFile).GetConfig(opts.SourceDir)

	scaff.Template(config, opts)
}
