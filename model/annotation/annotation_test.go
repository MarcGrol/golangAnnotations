package annotation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGarbage(t *testing.T) {
	_, ok := ResolveAnnotation(`// wvdwadbvb`)
	assert.False(t, ok)
}

func TestUnknownAction(t *testing.T) {
	RegisterAnnotation("Event", []string{}, validateAnnotation)

	_, ok := ResolveAnnotation(`// {"Annotation":"Haha","With":{"X":"Y"}}`)
	assert.False(t, ok)
}

func TestCorrectEventAnnotation(t *testing.T) {
	RegisterAnnotation("Event", []string{}, validateAnnotation)

	annotation, ok := ResolveAnnotation(`// {"Annotation":"Event","With":{"Aggregate":"Test"}}`)
	assert.True(t, ok)
	assert.Equal(t, "Event", annotation.Annotation)
	assert.Equal(t, "Test", annotation.With["Aggregate"])
}

func validateAnnotation(annot Annotation) bool {
	return true
}
