package main

import (
	"os"
	"path/filepath"

	"github.com/devshorts/scaff/scaff"
	"github.com/devshorts/scaff/scaff/file"
	"github.com/jessevdk/go-flags"
	"github.com/otiai10/copy"
	"github.com/sirupsen/logrus"
)

func main() {
	var opts struct {
		SourceDir       string `short:"d" long:"source_dir" description:"Source directory containing templates" required:"true"`
		TargetDir       string `short:"t" long:"target_dir" description:"Target directory to make with templated data" required:"true"`
		ScaffConfigFile string `long:"scaff_file" description:"Name of yaml file containing config. Defaults to .scaff.yml"`
		DryRun          bool   `long:"dry_run" description:"Dry Run"`
	}

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

	// given the key value pairs in the config, resolve them from stdin
	// and verify they are correct
	bag := scaff.NewBagResolver(os.Stdin, os.Stdout, config).ResolveBag()

	// create a new template runner
	templator := scaff.NewTemplator(config.FileConfig)

	// create the rules that allow for templating
	rules := scaff.NewRuleRunner(bag)

	copy.Copy(opts.SourceDir, opts.TargetDir)

	// template rules to the directories sorted from longest
	// to shortest. This allows us to rename and move directories
	// without having stale data in the path list
	for _, dir := range templator.GetAllDirs(opts.TargetDir) {
		templator.TemplatePath(dir, rules, opts.DryRun)
	}

	// once all directories have been moved apply the rules
	// to files and file names
	for _, file := range templator.GetAllFiles(opts.TargetDir) {
		templator.TemplateFile(file, rules, opts.DryRun)
	}

	if opts.ScaffConfigFile == "" {
		// clear out the .scaff file in the target
		os.Remove(filepath.Join(opts.TargetDir, scaff.DefaultSourcePath))
	}
}
