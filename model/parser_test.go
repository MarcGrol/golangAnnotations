package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFile(t *testing.T) {
	parsedSources, err := Parse("./gen_ast.json")
	assert.NoError(t, err)
	assert.NotEmpty(t, parsedSources.Structs)

	found := false
	for _, s := range parsedSources.Structs {
		if s.PackageName == "model" && s.Name == "ParsedSources" {
			found = true
			break
		}
	}
	assert.True(t, found)
}
