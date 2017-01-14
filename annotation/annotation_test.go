package annotation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGarbage(t *testing.T) {
	ClearRegisteredAnnotations()
	RegisterAnnotation("Event", []string{}, validateOk)

	assert.Empty(t, ResolveAnnotations([]string{`// wvdwadbvb`}))
}

func TestInvalidSyntax(t *testing.T) {
	ClearRegisteredAnnotations()
	RegisterAnnotation("Event", []string{}, validateOk)

	assert.Empty(t, ResolveAnnotations([]string{`// @X( a = "A" `}))
}

func TestTokensInValue(t *testing.T) {
	ClearRegisteredAnnotations()
	RegisterAnnotation("Event", []string{}, validateOk)

	ann, ok := ResolveAnnotationByName([]string{`// @SomethingElse( aggregate = "@A@")`, `// @Event( aggregate = "@A@")`}, "Event")
	assert.True(t, ok)
	assert.Equal(t, "@A@", ann.Attributes["aggregate"])
}

func TestUnknownName(t *testing.T) {
	ClearRegisteredAnnotations()
	RegisterAnnotation("X", []string{}, validateOk)

	assert.Empty(t, ResolveAnnotations([]string{`// @Y( a = "A" `}))
}

func TestCorrectAnnotation(t *testing.T) {
	ClearRegisteredAnnotations()
	RegisterAnnotation("X", []string{}, validateOk)

	ann, ok := ResolveAnnotationByName([]string{`// @X( a = "A" )`}, "X")
	assert.True(t, ok)
	assert.Equal(t, "X", ann.Name)
	assert.Equal(t, "A", ann.Attributes["a"])
}

func TestAnnotationWithValidationError(t *testing.T) {
	ClearRegisteredAnnotations()
	RegisterAnnotation("X", []string{}, validateError)

	assert.Empty(t, ResolveAnnotations([]string{`// @X( a = "A" )`}))
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
