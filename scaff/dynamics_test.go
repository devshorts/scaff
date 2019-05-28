package scaff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDynamics(t *testing.T) {
	config := NewParser("").GetConfig("test")

	result := ResolveDynamics(config, map[string]string{
		"foo": "foo",
		"biz": "bar",
	})

	assert.Equal(t, "hello/world/FOO/BAR", result["foo_biz"])
}
