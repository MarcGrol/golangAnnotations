package annotation

import (
	"testing"

	"log"

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
	RegisterAnnotation("Event", []string{}, validateOk)

	_, ok := ResolveAnnotation(`// {"Annotation":"Haha","With":{"X":"Y"}}`)
	assert.False(t, ok)
}

func TestCorrectEventAnnotation(t *testing.T) {
	ClearRegisteredAnnotations()
	RegisterAnnotation("Event", []string{}, validateOk)

	annotation, ok := ResolveAnnotation(`// {"Annotation":"Event","With":{"Aggregate":"Test"}}`)
	assert.True(t, ok)
	assert.Equal(t, "Event", annotation.Annotation)
	assert.Equal(t, "Test", annotation.With["Aggregate"])
}

func TestIcompleteEventAnnotation(t *testing.T) {
	ClearRegisteredAnnotations()
	RegisterAnnotation("Event", []string{}, validateError)

	_, ok := ResolveAnnotation(`// {"Annotation":"Event","With":{}}`)
	assert.False(t, ok)
}

func validateOk(annot Annotation) bool {
	log.Printf("good")
	return true
}

func validateError(annot Annotation) bool {
	log.Printf("wrong")
	return false
}
