package main

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilenameFiltering(t *testing.T) {
	var excludePattern = regexp.MustCompile(excludeMatchPattern)
	assert.False(t, excludePattern.MatchString("a.go"))
	assert.False(t, excludePattern.MatchString("a.txt"))
	assert.True(t, excludePattern.MatchString("gen_a.go"))
}
