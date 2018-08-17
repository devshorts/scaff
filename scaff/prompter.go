package scaff

import (
	"fmt"
	"bufio"
	"os"
)

type ContextPrompter struct{}

func NewPrompter() ContextPrompter {
	return ContextPrompter{}
}

func (c ContextPrompter) Resolve(config ScaffConfig) map[string]string {
	bag := make(map[string]string)

	reader := bufio.NewScanner(os.Stdin)

	for k, v := range config.Context {
		fmt.Print(string(v) + ": ")

		reader.Scan()

		bag[string(k)] = reader.Text()
	}

	return bag
}
