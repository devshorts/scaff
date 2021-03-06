package scaff

import (
	"testing"

	"github.com/devshorts/scaff/scaff/config"
	"github.com/stretchr/testify/assert"
)

func TestGetDirs(t *testing.T) {
	names := NewTemplator(config.FileConfig{}).GetAllDirs("../test")

	expectedDirs := []string{"../test/folder2/folder3", "../test/folder2", "../test/folder1", "../test"}

	assert.Equal(t, names, expectedDirs)
}

func TestPopSegment(t *testing.T) {
	resolver := NewTemplator(config.FileConfig{})

	path := "foo/bar/biz"
	segment, remaining := resolver.popSegment(path)
	assert.Equal(t, segment, "biz")
	assert.Equal(t, remaining, "foo/bar")

	path = "foo"
	segment, remaining = resolver.popSegment(path)
	assert.Equal(t, segment, "foo")
	assert.Equal(t, remaining, "")
}
