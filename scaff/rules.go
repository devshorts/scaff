package scaff

import (
	"regexp"
	"strings"
)

type ReplacementRule interface {
	Replace(text string) (string, bool)
}

type RuleName string

// Given a text text and a rule, see if the subsequent
// text matches the rule. For example "__camel_Foo__"
// should match "Foo" if the id is "camel".
func extractFormatToken(ruleName RuleName, text string) (string, bool) {
	re := regexp.MustCompile("__" + string(ruleName) + "_(.*)__")

	match := re.FindStringSubmatch(text)

	if len(match) < 2 {
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
		PackageRule{runner:runner},
		IdRule{runner:runner},
	}
}

func (runner RuleRunner) Replace(text string) string {
	for _, rule := range runner.getRules() {
		if replaced, ok := rule.Replace(text); ok {
			text = replaced
		}
	}

	return text
}

// Applies an id rule
func (runner RuleRunner) processText(
	text string,
	ruleName RuleName,
	processor func(string) string) (string, bool) {

	re := regexp.MustCompile("__" + string(ruleName) + "_(.*)__")

	result := re.ReplaceAllStringFunc(text, func(match string) string {
		if token, ok := extractFormatToken(ruleName, match); ok {
			if replace, ok := runner.ctx[token]; ok {
				return processor(replace)
			}
		}

		return match
	})


	return result, true
}

type CamelCase struct {
	runner RuleRunner
}

func (c CamelCase) Replace(text string) (string, bool) {
	return c.runner.processText(text, RuleName("camel"), func(s string) string {
		title := strings.Title(s)

		return strings.ToLower(string(title[0])) + title[1:]
	})
}

var _ ReplacementRule = CamelCase{}

type SnakeCase struct {
	runner RuleRunner
}

func (c SnakeCase) Replace(text string) (string, bool) {
	return c.runner.processText(text, RuleName("snake"), func(s string) string {
		return strings.Replace(s, " ", "_", -1)
	})
}

var _ ReplacementRule = SnakeCase{}

type UpperCase struct {
	runner RuleRunner
}

func (c UpperCase) Replace(text string) (string, bool) {
	return c.runner.processText(text, RuleName("upper"), func(s string) string {
		return strings.ToUpper(s)
	})
}

var _ ReplacementRule = UpperCase{}

type LowerCase struct {
	runner RuleRunner
}

func (c LowerCase) Replace(text string) (string, bool) {
	return c.runner.processText(text, RuleName("lower"), func(s string) string {
		return strings.ToLower(s)
	})
}

var _ ReplacementRule = LowerCase{}

// Replace text of a.b.c to a/b/c
type PackageRule struct {
	runner RuleRunner
}

func (c PackageRule) Replace(text string) (string, bool) {
	return c.runner.processText(text, RuleName("pkg"), func(s string) string {
		return strings.Replace(s, ".", "/", -1)
	})
}

var _ ReplacementRule = PackageRule{}

// Do nothing but swap placeholders
type IdRule struct {
	runner RuleRunner
}

func (c IdRule) Replace(text string) (string, bool) {
	return c.runner.processText(text, RuleName("id"), func(s string) string {
		return s
	})
}

var _ ReplacementRule = IdRule{}

