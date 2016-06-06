package annotation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGarbage(t *testing.T) {
	ClearRegisteredAnnotations()
	RegisterAnnotation("Event", []string{}, validateOk)

	_, ok := ResolveAnnotation(`// wvdwadbvb`)
	assert.False(t, ok)
}

func TestUnknownAction(t *testing.T) {
	ClearRegisteredAnnotations()
	RegisterAnnotation("X", []string{}, validateOk)

	_, ok := ResolveAnnotation(`// {"Annotation":"Y","With":{"a":"A"}}`)
	assert.False(t, ok)
}

func TestCorrectEventAnnotation(t *testing.T) {
	ClearRegisteredAnnotations()
	RegisterAnnotation("X", []string{}, validateOk)

	annotation, ok := ResolveAnnotation(`// {"Annotation":"X","With":{"a":"A"}}`)
	assert.True(t, ok)
	assert.Equal(t, "X", annotation.Annotation)
	assert.Equal(t, "A", annotation.With["a"])
}

func TestIcompleteEventAnnotation(t *testing.T) {
	ClearRegisteredAnnotations()
	RegisterAnnotation("X", []string{}, validateError)

	_, ok := ResolveAnnotation(`// {"Annotation":"X","With":{"a":"A"}}`)
	assert.False(t, ok)
}

func validateOk(annot Annotation) bool {
	return true
}

func validateError(annot Annotation) bool {
	return false
}
