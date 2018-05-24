package filegen

import (
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
