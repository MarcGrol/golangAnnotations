package annotation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGarbage(t *testing.T) {
	registry := NewRegistry([]AnnotationDescriptor{
		{
			Name:       "Event",
			ParamNames: []string{},
			Validator:  validateOk,
		},
	})

	assert.Empty(t, registry.ResolveAnnotations([]string{`// wvdwadbvb`}))
}

func TestInvalidSyntax(t *testing.T) {
	registry := NewRegistry([]AnnotationDescriptor{
		{
			Name:       "Event",
			ParamNames: []string{},
			Validator:  validateOk,
		},
	})

	assert.Empty(t, registry.ResolveAnnotations([]string{`// @X( a = "A" `}))
}

func TestTokensInValue(t *testing.T) {
	registry := NewRegistry([]AnnotationDescriptor{
		{
			Name:       "Event",
			ParamNames: []string{},
			Validator:  validateOk,
		},
	})

	ann, ok := registry.ResolveAnnotationByName([]string{`// @SomethingElse( aggregate = "@A@")`, `// @Event( aggregate = "@A@")`}, "Event")
	assert.True(t, ok)
	assert.Equal(t, "@A@", ann.Attributes["aggregate"])
}

func TestUnknownName(t *testing.T) {
	registry := NewRegistry([]AnnotationDescriptor{
		{
			Name:       "X",
			ParamNames: []string{},
			Validator:  validateOk,
		},
	})

	assert.Empty(t, registry.ResolveAnnotations([]string{`// @Y( a = "A" `}))
}

func TestCorrectAnnotation(t *testing.T) {
	registry := NewRegistry([]AnnotationDescriptor{
		{
			Name:       "X",
			ParamNames: []string{},
			Validator:  validateOk,
		},
	})

	ann, ok := registry.ResolveAnnotationByName([]string{`// @X( a = "A" )`}, "X")
	assert.True(t, ok)
	assert.Equal(t, "X", ann.Name)
	assert.Equal(t, "A", ann.Attributes["a"])
}

func TestAnnotationWithValidationError(t *testing.T) {
	registry := NewRegistry([]AnnotationDescriptor{
		{
			Name:       "X",
			ParamNames: []string{},
			Validator:  validateError,
		},
	})

	assert.Empty(t, registry.ResolveAnnotations([]string{`// @X( a = "A" )`}))
}

func TestAnnotationWithTypicalCharacters(t *testing.T) {
	registry := NewRegistry([]AnnotationDescriptor{
		{
			Name:       "Doit",
			ParamNames: []string{},
			Validator:  validateOk,
		},
	})

	annotation, ok := registry.ResolveAnnotation(`// @Doit( a="/A/", b="/B" )`)
	assert.True(t, ok)
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
