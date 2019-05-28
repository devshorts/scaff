package scaff

import (
	"testing"

	"github.com/devshorts/scaff/scaff/config"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	context := NewParser("").GetConfig("test").Context

	assert.Equal(t, context["biz"].Description, config.Description("baz"))
	assert.Equal(t, context["biz"].Default, "")
	assert.Equal(t, context["foo"].Description, config.Description("bar"))
	assert.Equal(t, context["foo"].Default, "default")

}
