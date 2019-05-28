package scaff

import (
	"os"
	"path/filepath"

	"github.com/devshorts/scaff/scaff/config"
	"github.com/otiai10/copy"
	"github.com/sirupsen/logrus"
)

type Opts struct {
	SourceDir       string `short:"d" long:"source_dir" description:"Source directory containing templates" required:"true"`
	TargetDir       string `short:"t" long:"target_dir" description:"Target directory to make with templated data" required:"true"`
	ScaffConfigFile string `long:"scaff_file" description:"Name of yaml file containing config. Defaults to .scaff.yml"`
	DryRun          bool   `long:"dry_run" description:"Dry Run"`
}

func Template(config config.ScaffConfig, opts Opts) {
	// given the key value pairs in the config, resolve them from stdin
	// and verify they are correct
	bag := NewBagResolver(os.Stdin, os.Stdout, config).ResolveBag()

	if err := TemplateWithConfig(config, bag, opts); err != nil {
		logrus.Fatal(err)
	}
}

func TemplateWithConfig(config config.ScaffConfig, bag map[string]string, opts Opts) error {
	// create a new template runner
	templator := NewTemplator(config.FileConfig)
	// create the rules that allow for templating
	rules := NewRuleRunner(bag)

	if err := copy.Copy(opts.SourceDir, opts.TargetDir); err != nil {
		return err
	}

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
		if err := os.Remove(filepath.Join(opts.TargetDir, DefaultSourcePath)); err != nil {
			return err
		}
	}

	return nil
}
