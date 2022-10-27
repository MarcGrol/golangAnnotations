package eventAnnotation

import (
	"testing"

	"github.com/f0rt/golangAnnotations/generator/annotation"
	"github.com/stretchr/testify/assert"
)

func TestCorrectEventAnnotation(t *testing.T) {
	registry := annotation.NewRegistry(Get())

	ann, ok := registry.ResolveAnnotationByName([]string{`// @Event( aggregate = "test", isRootEvent = "true" )`}, "Event")
	assert.True(t, ok)
	assert.Equal(t, "test", ann.Attributes["aggregate"])
	assert.Equal(t, "true", ann.Attributes["isrootevent"])
}

func TestIncompleteEventAnnotation(t *testing.T) {
	registry := annotation.NewRegistry(Get())

	assert.Empty(t, registry.ResolveAnnotations([]string{`// @Event()`}))
}

func TestEmptyEventAnnotation(t *testing.T) {
	registry := annotation.NewRegistry(Get())

	assert.Empty(t, registry.ResolveAnnotations([]string{`// @Event( aggregate = "")`}))
}
