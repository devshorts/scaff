package lang

import (
	"strings"

	"github.com/devshorts/scaff/scaff/config"
)

type LangProcessor interface {
	Process(contents string) string
}

type GoProcessor struct {
	config          config.GoRules
	replacedPackage string
}

func NewGoProcessor(config config.GoRules, context map[string]string) *GoProcessor {
	return &GoProcessor{config: config, replacedPackage: context[config.ReplaceRule]}
}

func (g GoProcessor) Process(contents string) string {
	return strings.ReplaceAll(contents, g.config.SourcePackage, g.replacedPackage)
}

var _ LangProcessor = GoProcessor{}
