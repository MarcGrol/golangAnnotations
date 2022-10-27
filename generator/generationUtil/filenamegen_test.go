package generationUtil

import (
	"testing"

	"github.com/MarcGrol/golangAnnotations/generator"
	"github.com/stretchr/testify/assert"
)

func TestNoDir(t *testing.T) {
	assert.Equal(t, generator.GenfilePrefix+"test.txt", Prefixed("test.txt"))
}

func TestWithRelativeDir(t *testing.T) {
	assert.Equal(t, "dir/"+generator.GenfilePrefix+"test.txt", Prefixed("dir/test.txt"))
}

func TestWithAbsoluteDir(t *testing.T) {
	assert.Equal(t, "/dir/"+generator.GenfilePrefix+"test.txt", Prefixed("/dir/test.txt"))
}
