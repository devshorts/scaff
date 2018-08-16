package scaff

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	context := NewParser("").GetConfig("test").Context

	assert.Equal(t, context["biz"], Description("baz"))
	assert.Equal(t, context["foo"], Description("bar"))
}
