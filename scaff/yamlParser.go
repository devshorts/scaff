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

func (p Parser) GetConfig(path string) ScaffConfig {
	bytes, _ := ioutil.ReadFile(filepath.Join(path, p.configFileName))

	var config ScaffConfig

	yaml.Unmarshal(bytes, &config)

	return config
}
