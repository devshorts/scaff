package scaff

import (
	"fmt"
	"bufio"
	"io"
)

type ContextPrompter struct{}

func NewPrompter() ContextPrompter {
	return ContextPrompter{}
}

// Asks the user to supply the results
func (c ContextPrompter) ResolveBag(config ScaffConfig, reader io.Reader) map[string]string {
	bag := make(map[string]string)

	scanner := bufio.NewScanner(reader)

	for k, v := range config.Context {
		defaultDescription := ""

		if len(v.Default) > 0 {
			defaultDescription = " (" + v.Default + ")"
		}

		fmt.Print(string(v.Description) + defaultDescription + ": ")

		scanner.Scan()

		result := scanner.Text()

		if len(result) == 0 && len(v.Default) > 0 {
			result = v.Default
		}

		bag[string(k)] = result
	}

	return bag
}

func (c ContextPrompter) ConfirmBag(bag map[string]string, config ScaffConfig, writer io.Writer, reader io.Reader) {
	fmt.Fprintln(writer)
	fmt.Fprintln(writer, "!Verify!")
	fmt.Fprintln(writer)

	for k, v := range bag {
		desc := config.Context[Name(k)].Description

		fmt.Fprintln(writer, fmt.Sprintf("%s = %s", desc, v))
	}

	fmt.Fprintln(writer)
	fmt.Fprint(writer, "Confirm?...")

	bufio.NewScanner(reader).Scan()
}

