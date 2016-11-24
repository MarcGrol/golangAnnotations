package eventAnnotation

import (
	"github.com/MarcGrol/golangAnnotations/annotation"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCorrectEventAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	annot, ok := annotation.ResolveAnnotations([]string{`// @Event( aggregate = "test", isRootEvent = "true" )`})
	assert.True(t, ok)
	assert.Equal(t, "test", annot.Attributes["aggregate"])
	assert.Equal(t, "true", annot.Attributes["isrootevent"])
}

func TestIncompleteEventAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	_, ok := annotation.ResolveAnnotations([]string{`// @Event()`})
	assert.False(t, ok)
}

func TestEmptyEventAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	_, ok := annotation.ResolveAnnotations([]string{`// @Event( aggregate = "")`})
	assert.False(t, ok)
}
