package filegen

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoDir(t *testing.T) {
	assert.Equal(t, genfilePrefix+"test.txt", Prefixed("test.txt"))
}

func TestWithRelativeDir(t *testing.T) {
	assert.Equal(t, "dir/"+genfilePrefix+"test.txt", Prefixed("dir/test.txt"))
}

func TestWithAbsoluteDir(t *testing.T) {
	assert.Equal(t, "/dir/"+genfilePrefix+"test.txt", Prefixed("/dir/test.txt"))
}

func TestFilenameMatch(t *testing.T) {
	var excludePattern = regexp.MustCompile(ExcludeMatchPattern())
	assert.False(t, excludePattern.MatchString("a.go"))
	assert.False(t, excludePattern.MatchString("a.txt"))
	assert.True(t, excludePattern.MatchString("gen_a.go"))
}
