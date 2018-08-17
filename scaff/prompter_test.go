package scaff

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestPrompter(t *testing.T) {
	t.Skip()

	c := ScaffConfig{
		Context: map[Name]Description{
			Name("test"): Description("Give me your tests"),
		},
	}

	result := NewPrompter().ResolveBag(c)

	assert.Equal(t, result["test"], "foo")
}
