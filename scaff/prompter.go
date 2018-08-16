package scaff

import (
	"fmt"
)

type ContextPrompter struct{}

func NewPrompter() ContextPrompter {
	return ContextPrompter{}
}

func (c ContextPrompter) Resolve(config ScaffConfig) map[string]string {
	bag := make(map[string]string)

	for k, v := range config.Context {
		fmt.Print(string(v) + ": ")

		var input string

		fmt.Scanln(&input)

		bag[string(k)] = input
	}

	return bag
}
