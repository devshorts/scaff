package scaff

import (
	"fmt"
	"bufio"
	"io"
	"github.com/devshorts/scaff/scaff/sstring"
	"os/exec"
	"os"
	"github.com/sirupsen/logrus"
)

type BagResolver struct {
	stdin  io.Reader
	out    io.Writer
	config ScaffConfig
}

func NewBagResolver(stdin io.Reader, out io.Writer, config ScaffConfig) BagResolver {
	return BagResolver{
		stdin:  stdin,
		out:    out,
		config: config,
	}
}

// Asks the user to supply the results
func (c BagResolver) ResolveBag() map[string]string {
	bag := c.parseBag()

	c.confirmBag(bag)

	c.postHookVerify(bag)

	return bag.AsRaw()
}

func (c BagResolver) parseBag() ResolvedConfig {
	bag := make(ResolvedConfig)

	scanner := bufio.NewScanner(c.stdin)

	for k, v := range c.config.Context {
		defaultDescription := ""

		if len(v.Default) > 0 {
			defaultDescription = " (" + v.Default + ")"
		}

		result := c.parseKeyFromInput(v, defaultDescription, scanner)

		bag[k] = ParsedValue{
			Source:      v,
			ParsedValue: result,
		}
	}

	return bag
}

func (c BagResolver) parseKeyFromInput(v TemplateValue, defaultDescription string, scanner *bufio.Scanner) string {
	fmt.Fprint(c.out, string(v.Description)+defaultDescription+": ")

	scanner.Scan()

	result := scanner.Text()

	if len(result) == 0 && len(v.Default) > 0 {
		result = v.Default
	}

	if sstring.IsEmpty(result) {
		fmt.Fprintln(c.out, "Please set this field")

		result = c.parseKeyFromInput(v, defaultDescription, scanner)
	}

	return result
}

func (c BagResolver) confirmBag(bag ResolvedConfig) {
	fmt.Fprintln(c.out)
	fmt.Fprintln(c.out, "!Verify!")
	fmt.Fprintln(c.out)

	for k, v := range bag {
		desc := c.config.Context[TemplateKey(k)].Description

		fmt.Fprintln(c.out, fmt.Sprintf("%s = %s", desc, v.ParsedValue))
	}

	fmt.Fprintln(c.out)
	fmt.Fprint(c.out, "Confirm?...")

	bufio.NewScanner(c.stdin).Scan()
}

func (c BagResolver) postHookVerify(configs ResolvedConfig) {
	for _, v := range configs {
		if !sstring.IsEmpty(v.Source.VerifyHook.Command) {
			args := v.Source.VerifyHook.Args

			if !sstring.IsEmpty(v.ParsedValue) {
				args = append(args, v.ParsedValue)
			}

			cmd := exec.Command(v.Source.VerifyHook.Command, args...)

			logrus.Info(fmt.Sprintf("Verifying %s with '%s %s'",
				v.Source.Description,
				v.Source.VerifyHook.Command,
				args))

			err := cmd.Run()

			if err != nil {
				logrus.Error(fmt.Sprintf("Error %s", err))

				os.Exit(-1)
			}
		}
	}
}
