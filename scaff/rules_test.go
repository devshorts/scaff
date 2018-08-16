package scaff

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCamelCase_Replace(t *testing.T) {
	bag := map[string]string {
		"test": "boo burns",
	}

	formatter := NewRuleFormatter(bag)

	assert.Equal(t, formatter.Replace("__camel_test__"), "boo Burns")
}

func TestLowerCase_Replace(t *testing.T) {
	bag := map[string]string {
		"test": "BOOURNS",
	}

	formatter := NewRuleFormatter(bag)

	assert.Equal(t, formatter.Replace("__lower_test__"), "boourns")
}

func TestSnakeCase_Replace(t *testing.T) {
	bag := map[string]string {
		"test": "boo urns",
	}

	formatter := NewRuleFormatter(bag)

	assert.Equal(t, formatter.Replace("__snake_test__"), "boo_urns")
}

func TestUpperCase_Replace(t *testing.T) {
	bag := map[string]string {
		"test": "boo urns",
	}

	formatter := NewRuleFormatter(bag)

	assert.Equal(t, formatter.Replace("__upper_test__"), "BOO URNS")
}
