package scaff

import (
	"regexp"
	"strings"
)

type ReplacementRule interface {
	Replace(text string) (string, bool)
}

type RuleName string

func (RuleRunner) tokenRegex(ruleName RuleName, tokenDelimiter string) *regexp.Regexp {
	// if there are special chars in the token, escape them
	escapedRegexValues := []string{"$"}

	for _, escapedValue := range escapedRegexValues {
		tokenDelimiter = strings.Replace(tokenDelimiter, escapedValue, "\\"+escapedValue, -1)
	}

	re := regexp.MustCompile(tokenDelimiter + string(ruleName) + "_(.*)" + tokenDelimiter)

	return re
}

// Given a text text and a rule, see if the subsequent
// text matches the rule. For example "__camel_Foo__"
// should match "Foo" if the id is "camel".
func (runner RuleRunner) extractFormatToken(ruleName RuleName, text string, tokenDelimiter string) (string, bool) {
	re := runner.tokenRegex(ruleName, tokenDelimiter)

	match := re.FindStringSubmatch(text)

	if len(match) < 2 {
		return "", false
	}

	return match[1], true
}

type RuleRunner struct {
	ctx map[string]string
}

func NewRuleRunner(bag map[string]string) RuleRunner {
	return RuleRunner{
		ctx: bag,
	}
}

func (runner RuleRunner) getRules(tokenDelimiter string) []ReplacementRule {
	return []ReplacementRule{
		CamelCase{TokenRule{runner: runner, tokenDelimiter: tokenDelimiter}},
		SnakeCase{TokenRule{runner: runner, tokenDelimiter: tokenDelimiter}},
		LowerCase{TokenRule{runner: runner, tokenDelimiter: tokenDelimiter}},
		UpperCase{TokenRule{runner: runner, tokenDelimiter: tokenDelimiter}},
		PackageRule{TokenRule{runner: runner, tokenDelimiter: tokenDelimiter}},
		IdRule{TokenRule{runner: runner, tokenDelimiter: tokenDelimiter}},
	}
}

func (runner RuleRunner) Replace(text string, tokenDelimiter string) string {
	for _, rule := range runner.getRules(tokenDelimiter) {
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
	tokenDelimiter string,
	processor func(string) string) (string, bool) {

	re := runner.tokenRegex(ruleName, tokenDelimiter)

	result := re.ReplaceAllStringFunc(text, func(match string) string {
		if token, ok := runner.extractFormatToken(ruleName, match, tokenDelimiter); ok {
			if replace, ok := runner.ctx[token]; ok {
				return processor(replace)
			}
		}

		return match
	})

	return result, true
}

type TokenRule struct {
	runner         RuleRunner
	tokenDelimiter string
}

type CamelCase struct {
	TokenRule
}

func (c CamelCase) Replace(text string) (string, bool) {
	return c.runner.processText(
		text,
		RuleName("camel"),
		c.tokenDelimiter,
		func(s string) string {
			title := strings.Title(s)

			return strings.ToLower(string(title[0])) + title[1:]
		})
}

var _ ReplacementRule = CamelCase{}

type SnakeCase struct {
	TokenRule
}

func (c SnakeCase) Replace(text string) (string, bool) {
	return c.runner.processText(
		text,
		RuleName("snake"),
		c.tokenDelimiter,
		func(s string) string {
			return strings.Replace(s, " ", "_", -1)
		})
}

var _ ReplacementRule = SnakeCase{}

type UpperCase struct {
	TokenRule
}

func (c UpperCase) Replace(text string) (string, bool) {
	return c.runner.processText(
		text,
		RuleName("upper"),
		c.tokenDelimiter,
		func(s string) string {
			return strings.ToUpper(s)
		})
}

var _ ReplacementRule = UpperCase{}

type LowerCase struct {
	TokenRule
}

func (c LowerCase) Replace(text string) (string, bool) {
	return c.runner.processText(
		text,
		RuleName("lower"),
		c.tokenDelimiter,
		func(s string) string {
			return strings.ToLower(s)
		})
}

var _ ReplacementRule = LowerCase{}

// Replace text of a.b.c to a/b/c
type PackageRule struct {
	TokenRule
}

func (c PackageRule) Replace(text string) (string, bool) {
	return c.runner.processText(
		text,
		RuleName("pkg"),
		c.tokenDelimiter,
		func(s string) string {
			return strings.Replace(s, ".", "/", -1)
		})
}

var _ ReplacementRule = PackageRule{}

// Do nothing but swap placeholders
type IdRule struct {
	TokenRule
}

func (c IdRule) Replace(text string) (string, bool) {
	return c.runner.processText(
		text,
		RuleName("id"),
		c.tokenDelimiter,
		func(s string) string {
			return s
		})
}

var _ ReplacementRule = IdRule{}
