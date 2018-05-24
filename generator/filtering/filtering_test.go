package filtering

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilenameMatch(t *testing.T) {
	var excludePattern = regexp.MustCompile(ExcludeMatchPattern())
	assert.False(t, excludePattern.MatchString("a.go"))
	assert.False(t, excludePattern.MatchString("a.txt"))
	assert.True(t, excludePattern.MatchString("gen_a.go"))
}
