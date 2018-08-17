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

	prompter := scaff.NewPrompter()

	bag := prompter.ResolveBag(config, os.Stdin)

	prompter.ConfirmBag(bag, config, os.Stdout, os.Stdin)

	templator := scaff.NewTemplator()

	rules := scaff.NewRuleFormatter(bag)

	for _, dir := range templator.GetAllDirs(opts.Dir) {
		templator.TemplatePath(dir, rules, opts.DryRun)
	}

	for _, file := range templator.GetAllFiles(opts.Dir) {
		templator.TemplateFile(file, rules, opts.DryRun)
	}
}
