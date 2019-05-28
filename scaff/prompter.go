package scaff

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"

	"github.com/devshorts/scaff/scaff/config"
	"github.com/devshorts/scaff/scaff/sstring"
	"github.com/sirupsen/logrus"
)

type BagResolver struct {
	stdin  io.Reader
	out    io.Writer
	config config.ScaffConfig
}

func NewBagResolver(stdin io.Reader, out io.Writer, config config.ScaffConfig) BagResolver {
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

	raw := bag.AsRaw()

	return raw
}

func (c BagResolver) parseBag() config.ResolvedConfig {
	bag := make(config.ResolvedConfig)

	scanner := bufio.NewScanner(c.stdin)

	for k, v := range c.config.Context {
		defaultDescription := ""

		if len(v.Default) > 0 {
			defaultDescription = " (" + v.Default + ")"
		}

		bag[k] = c.parseKeyFromInput(v, defaultDescription, scanner)
	}

	return bag
}

func (c BagResolver) parseKeyFromInput(v config.TemplateValue, defaultDescription string, scanner *bufio.Scanner) config.ParsedValue {
	fmt.Fprint(c.out, string(v.Description)+defaultDescription+": ")

	scanner.Scan()

	userInput := scanner.Text()

	if len(userInput) == 0 && len(v.Default) > 0 {
		userInput = v.Default
	}

	result := config.ParsedValue{
		Source:      v,
		ParsedValue: userInput,
	}

	if sstring.IsEmpty(result.ParsedValue) || !c.postHookVerify(result) {
		fmt.Fprintln(c.out, fmt.Sprintf("A value of '%s' is invalid, please set it again", userInput))

		result = c.parseKeyFromInput(v, defaultDescription, scanner)
	}

	return result
}

func (c BagResolver) confirmBag(bag config.ResolvedConfig) {
	fmt.Fprintln(c.out)
	fmt.Fprintln(c.out, "!Verify!")
	fmt.Fprintln(c.out)

	for k, v := range bag {
		desc := c.config.Context[config.TemplateKey(k)].Description

		fmt.Fprintln(c.out, fmt.Sprintf("%s = %s", desc, v.ParsedValue))
	}

	fmt.Fprintln(c.out)
	fmt.Fprint(c.out, "Confirm?...")

	bufio.NewScanner(c.stdin).Scan()
}

func (c BagResolver) postHookVerify(parsed config.ParsedValue) bool {
	if !sstring.IsEmpty(parsed.Source.VerifyHook.Command) {
		args := parsed.Source.VerifyHook.Args

		if !sstring.IsEmpty(parsed.ParsedValue) {
			args = append(args, parsed.ParsedValue)
		}

		cmd := exec.Command(parsed.Source.VerifyHook.Command, args...)

		logrus.Debug(fmt.Sprintf("Verifying %s with '%s %s'",
			parsed.Source.Description,
			parsed.Source.VerifyHook.Command,
			args))

		err := cmd.Run()

		if err != nil {
			logrus.Warn(fmt.Sprintf("Error %s", err))

			return false
		}
	}

	return true
}
