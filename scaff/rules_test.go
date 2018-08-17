package scaff

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCamelCase_Replace(t *testing.T) {
	bag := map[string]string{
		"test": "boo burns",
	}

	formatter := NewRuleRunner(bag)

	assert.Equal(t, formatter.Replace("__camel_test__", DEFAULT_DELIM), "boo Burns")
}

func TestCamelCaseFullText_Replace(t *testing.T) {
	bag := map[string]string{
		"test": "boo burns",
	}

	formatter := NewRuleRunner(bag)

	assert.Equal(t, `
foo bar boo Burns test foo
bar foo boo Burns bar foo
`, formatter.Replace(`
foo bar __camel_test__ test foo
bar foo __camel_test__ bar foo
`, DEFAULT_DELIM))
}

func TestLowerCase_Replace(t *testing.T) {
	bag := map[string]string{
		"test": "BOOURNS",
	}

	formatter := NewRuleRunner(bag)

	assert.Equal(t, formatter.Replace("__lower_test__", DEFAULT_DELIM), "boourns")
}

func TestIdRule(t *testing.T) {
	bag := map[string]string{
		"test": "BOOURNS",
	}

	formatter := NewRuleRunner(bag)

	assert.Equal(t, formatter.Replace("__id_test__", DEFAULT_DELIM), "BOOURNS")
}

func TestPkgRule(t *testing.T) {
	bag := map[string]string{
		"test": "a.b.c",
	}

	formatter := NewRuleRunner(bag)

	assert.Equal(t, formatter.Replace("__pkg_test__", DEFAULT_DELIM), "a/b/c")
}

func TestSnakeCase_Replace(t *testing.T) {
	bag := map[string]string{
		"test": "boo urns",
	}

	formatter := NewRuleRunner(bag)

	assert.Equal(t, formatter.Replace("__snake_test__", DEFAULT_DELIM), "boo_urns")
}

func TestUpperCase_Replace(t *testing.T) {
	bag := map[string]string{
		"test": "boo urns",
	}

	formatter := NewRuleRunner(bag)

	assert.Equal(t, formatter.Replace("__upper_test__", DEFAULT_DELIM), "BOO URNS")
}

func TestCustomDelim(t *testing.T) {
	bag := map[string]string{
		"test": "boo urns",
	}

	formatter := NewRuleRunner(bag)

	assert.Equal(t, formatter.Replace("$$upper_test$$", "$$"), "BOO URNS")
}