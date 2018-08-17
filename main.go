package main

import (
	"github.com/jessevdk/go-flags"
	"os"
	"github.com/devshorts/scaff/scaff"
)

func main() {
	var opts struct {
		Dir string `short:"d" long:"directory" description:"Source directory" required:"true"`
		DryRun bool `long:"dry_run" description:"Dry Run"`
	}

	parser := flags.NewParser(&opts, flags.Default)

	if _, e := parser.Parse(); e != nil {
		os.Exit(-1)
	}

	config := scaff.NewParser("").GetConfig(opts.Dir)

	bag := scaff.NewPrompter().Resolve(config)

	resolver := scaff.NewFileResolver()

	rules := scaff.NewRuleFormatter(bag)

	for _, dir := range resolver.GetAllDirs(opts.Dir) {
		resolver.TemplatePath(dir, rules, opts.DryRun)
	}

	for _, file := range resolver.GetAllFiles(opts.Dir) {
		resolver.TemplateFile(file, rules, opts.DryRun)
	}
}
