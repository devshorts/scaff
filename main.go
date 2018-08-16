package main

import (
	"flag"
	"github.com/jessevdk/go-flags"
	"os"
	"github.com/devshorts/scaff/scaff"
)

func main() {
	var opts struct {
		Dir string `short:"d" long:"directory" description:"Source directory" required:"true"`
	}

	if _, e := flags.ParseArgs(&opts, flag.Args()); e != nil {
		os.Exit(-1)
	}

	config := scaff.NewParser("").GetConfig(opts.Dir)

	bag := scaff.NewPrompter().Resolve(config)

	resolver := scaff.NewFileResolver()

	rules := scaff.NewRuleFormatter(bag)


}
