package jsonAnnotation

import (
	"testing"

	"github.com/MarcGrol/golangAnnotations/annotation"
	"github.com/stretchr/testify/assert"
)

func TestCorrectEventAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	assert.NotEmpty(t, annotation.ResolveAnnotations([]string{`// @JsonEnum( )`}))
}

func TestEmptyEventAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	assert.Empty(t, annotation.ResolveAnnotations([]string{``}))
}
