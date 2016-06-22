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

func TestInvalidSyntax(t *testing.T) {
	ClearRegisteredAnnotations()
	RegisterAnnotation("Event", []string{}, validateOk)

	_, ok := ResolveAnnotation(`// @X( a = "A" `)
	assert.False(t, ok)
}

func TestTokensInValue(t *testing.T) {
	ClearRegisteredAnnotations()
	RegisterAnnotation("Event", []string{}, validateOk)

	annotation, ok := ResolveAnnotation(`// @Event( aggregate = "@A@")`)
	assert.True(t, ok)
	assert.Equal(t, "@A@", annotation.Attributes["aggregate"])
}

func TestUnknownName(t *testing.T) {
	ClearRegisteredAnnotations()
	RegisterAnnotation("X", []string{}, validateOk)

	_, ok := ResolveAnnotation(`// @Y( a = "A" `)
	assert.False(t, ok)
}

func TestCorrectAnnotation(t *testing.T) {
	ClearRegisteredAnnotations()
	RegisterAnnotation("X", []string{}, validateOk)

	annotation, ok := ResolveAnnotation(`// @X( a = "A" )`)
	assert.True(t, ok)
	assert.Equal(t, "X", annotation.Name)
	assert.Equal(t, "A", annotation.Attributes["a"])
}

func TestAnnotationWithValidationError(t *testing.T) {
	ClearRegisteredAnnotations()
	RegisterAnnotation("X", []string{}, validateError)

	_, ok := ResolveAnnotation(`// @X( a = "A" )`)
	assert.False(t, ok)
}

func TestAnnotationWithTypicalCharacters(t *testing.T) {
	ClearRegisteredAnnotations()
	RegisterAnnotation("Doit", []string{}, validateOk)

	annotation, err := parseAnnotation(`// @Doit( a="/A/", b="/B" )`)
	assert.NoError(t, err)
	assert.Equal(t, "Doit", annotation.Name)
	assert.Equal(t, "/A/", annotation.Attributes["a"])
	assert.Equal(t, "/B", annotation.Attributes["b"])
}

func validateOk(annot Annotation) bool {
	return true
}

func validateError(annot Annotation) bool {
	return false
}
