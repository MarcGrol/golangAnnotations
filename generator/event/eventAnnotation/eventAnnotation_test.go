package eventAnnotation

import (
	"testing"

	"github.com/MarcGrol/golangAnnotations/annotation"
	"github.com/stretchr/testify/assert"
)

func TestCorrectEventAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	ann, ok := annotation.ResolveAnnotationByName([]string{`// @Event( aggregate = "test", isRootEvent = "true" )`}, "Event")
	assert.True(t, ok)
	assert.Equal(t, "test", ann.Attributes["aggregate"])
	assert.Equal(t, "true", ann.Attributes["isrootevent"])
}

func TestIncompleteEventAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	assert.Empty(t, annotation.ResolveAnnotations([]string{`// @Event()`}))
}

func TestEmptyEventAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	assert.Empty(t, annotation.ResolveAnnotations([]string{`// @Event( aggregate = "")`}))
}
