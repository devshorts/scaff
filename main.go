package main

import (
	"github.com/jessevdk/go-flags"
	"os"
	"github.com/devshorts/scaff/scaff"
)

func main() {
	var opts struct {
		Dir    string `short:"d" long:"directory" description:"Source directory" required:"true"`
		ScaffConfigFile string `long:"scaff_file" description:"Name of yaml file containing config. Defaults to .scaff.yml"`
		DryRun bool   `long:"dry_run" description:"Dry Run"`
	}

	parser := flags.NewParser(&opts, flags.Default)

	if _, e := parser.Parse(); e != nil {
		os.Exit(-1)
	}

	// load a config from the scaff config file in the template directory
	config := scaff.NewParser(opts.ScaffConfigFile).GetConfig(opts.Dir)

	// given the key value pairs in the config, resolve them from stdin
	// and verify they are correct
	bag := scaff.NewBagResolver(os.Stdin, os.Stdout, config).ResolveBag()

	// create a new template runner
	templator := scaff.NewTemplator(config.FileConfig)

	// create the rules that allow for templating
	rules := scaff.NewRuleRunner(bag)

	// template rules to the directories sorted from longest
	// to shortest. This allows us to rename and move directories
	// without having stale data in the path list
	for _, dir := range templator.GetAllDirs(opts.Dir) {
		templator.TemplatePath(dir, rules, opts.DryRun)
	}

	// once all directories have been moved apply the rules
	// to files and file names
	for _, file := range templator.GetAllFiles(opts.Dir) {
		templator.TemplateFile(file, rules, opts.DryRun)
	}
}
