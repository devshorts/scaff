package scaff

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

const defaultSourcePath string = ".scaff.yml"

type Parser struct {
	configFileName string
}

func NewParser(sourceFile string) Parser {
	target := defaultSourcePath

	if sourceFile != "" {
		target = sourceFile
	}

	return Parser{configFileName: target}
}

type Name string
type Description string

type TemplateValue struct {
	Default string
	Description Description
}

type ScaffConfig struct {
	Context map[Name]TemplateValue
}

func (p Parser) GetConfig(path string) ScaffConfig {
	bytes, _ := ioutil.ReadFile(filepath.Join(path, p.configFileName))

	var config ScaffConfig

	yaml.Unmarshal(bytes, &config)

	return config
}
