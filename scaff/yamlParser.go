package scaff

import (
	"io/ioutil"
	"path/filepath"

	"github.com/devshorts/scaff/scaff/config"
	"gopkg.in/yaml.v2"
)

const DefaultSourcePath string = ".scaff.yml"

type Parser struct {
	configFileName string
}

func NewParser(sourceFile string) Parser {
	target := DefaultSourcePath

	if sourceFile != "" {
		target = sourceFile
	}

	return Parser{configFileName: target}
}

func (p Parser) GetConfig(path string) config.ScaffConfig {
	bytes, _ := ioutil.ReadFile(filepath.Join(path, p.configFileName))

	var config config.ScaffConfig

	yaml.Unmarshal(bytes, &config)

	return config
}
