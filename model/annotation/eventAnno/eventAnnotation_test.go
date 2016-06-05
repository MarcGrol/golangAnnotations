package eventAnno

import (
	"testing"

	"github.com/MarcGrol/astTools/model/annotation"
	"github.com/stretchr/testify/assert"
)

func TestCorrectEventAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	annot, ok := annotation.ResolveAnnotations([]string{`// {"Annotation":"Event","With":{"Aggregate":"test"}}`})
	assert.True(t, ok)
	assert.Equal(t, "test", annot.With["Aggregate"])
}

func TestIncompleteEventAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	_, ok := annotation.ResolveAnnotations([]string{`// {"Annotation":"Event"}`})
	assert.False(t, ok)
}

func TestEmptyEventAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	_, ok := annotation.ResolveAnnotations([]string{`// {"Annotation":"Event""With":{"Aggregate":""}}`})
	assert.False(t, ok)
}
