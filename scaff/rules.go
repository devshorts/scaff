package scaff

import (
	"regexp"
	"strings"
)

type ReplacementRule interface {
	Replace(segment string) (string, bool)
}

type RuleName string

// Given a text segment and a rule, see if the subsequent
// text matches the rule. For example "__camel_Foo__"
// should match "Foo" if the id is "camel".
func extractFormatToken(ruleName RuleName, segment string) (string, bool) {
	re := regexp.MustCompile("__" + string(ruleName) + "_(.*)__")

	match := re.FindStringSubmatch(segment)

	if len(match) != 2 {
		return "", false
	}

	return match[1], true
}

type RuleRunner struct {
	ctx map[string]string
}

func NewRuleFormatter(bag map[string]string) RuleRunner {
	return RuleRunner{ctx: bag}
}

func (runner RuleRunner) getRules() []ReplacementRule {
	return []ReplacementRule{
		CamelCase{runner: runner},
		SnakeCase{runner: runner},
		LowerCase{runner: runner},
		UpperCase{runner: runner},
	}
}

func (runner RuleRunner) Replace(segment string) string {
	for _, rule := range runner.getRules() {
		if replaced, ok := rule.Replace(segment); ok {
			return replaced
		}
	}

	return segment
}

// Applies an id rule
func (runner RuleRunner) processSegment(
	segment string,
	ruleName RuleName,
	processor func(string) string) (string, bool) {
	if result, ok := extractFormatToken(ruleName, segment); ok {
		if replace, ok := runner.ctx[result]; ok {
			return processor(replace), true
		}
	}

	return "", false
}

type CamelCase struct {
	runner RuleRunner
}

func (c CamelCase) Replace(segment string) (string, bool) {
	return c.runner.processSegment(segment, RuleName("camel"), func(s string) string {
		title := strings.Title(s)

		return strings.ToLower(string(title[0])) + title[1:]
	})
}

var _ ReplacementRule = CamelCase{}

type SnakeCase struct {
	runner RuleRunner
}

func (c SnakeCase) Replace(segment string) (string, bool) {
	return c.runner.processSegment(segment, RuleName("snake"), func(s string) string {
		return strings.Replace(s, " ", "_", -1)
	})
}

var _ ReplacementRule = SnakeCase{}

type UpperCase struct {
	runner RuleRunner
}

func (c UpperCase) Replace(segment string) (string, bool) {
	return c.runner.processSegment(segment, RuleName("upper"), func(s string) string {
		return strings.ToUpper(s)
	})
}

var _ ReplacementRule = UpperCase{}

type LowerCase struct {
	runner RuleRunner
}

func (c LowerCase) Replace(segment string) (string, bool) {
	return c.runner.processSegment(segment, RuleName("lower"), func(s string) string {
		return strings.ToLower(s)
	})
}

var _ ReplacementRule = LowerCase{}
