package jsonAnnotation

import (
	"testing"

	"github.com/MarcGrol/golangAnnotations/annotation"
	"github.com/stretchr/testify/assert"
)

func TestCorrectEventAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	_, ok := annotation.ResolveAnnotations([]string{`// @JsonEnum( )`})
	assert.True(t, ok)
}

func TestEmptyEventAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	_, ok := annotation.ResolveAnnotations([]string{``})
	assert.False(t, ok)
}
