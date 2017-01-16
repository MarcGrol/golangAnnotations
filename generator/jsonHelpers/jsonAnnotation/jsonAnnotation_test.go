package jsonAnnotation

import (
	"testing"

	"github.com/MarcGrol/golangAnnotations/annotation"
	"github.com/stretchr/testify/assert"
)

func TestCorrectEnumAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	assert.NotEmpty(t, annotation.ResolveAnnotations([]string{`// @JsonEnum( )`}))
}

func TestEmptyEnumAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	assert.Empty(t, annotation.ResolveAnnotations([]string{``}))
}

func TestCorrectStructAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	assert.NotEmpty(t, annotation.ResolveAnnotations([]string{`// @JsonStruct( )`}))
}

func TestEmptyStructAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	assert.Empty(t, annotation.ResolveAnnotations([]string{``}))
}
