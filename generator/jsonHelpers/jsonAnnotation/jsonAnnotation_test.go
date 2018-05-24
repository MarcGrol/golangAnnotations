package jsonAnnotation

import (
	"testing"

	"github.com/MarcGrol/golangAnnotations/generator/annotation"
	"github.com/stretchr/testify/assert"
)

func TestCorrectEnumAnnotation(t *testing.T) {
	registry := annotation.NewRegistry(Get())

	assert.NotEmpty(t, registry.ResolveAnnotations([]string{`// @JsonEnum( )`}))
}

func TestEmptyEnumAnnotation(t *testing.T) {
	registry := annotation.NewRegistry(Get())

	assert.Empty(t, registry.ResolveAnnotations([]string{``}))
}

func TestCorrectStructAnnotation(t *testing.T) {
	registry := annotation.NewRegistry(Get())

	assert.NotEmpty(t, registry.ResolveAnnotations([]string{`// @JsonStruct( )`}))
}

func TestEmptyStructAnnotation(t *testing.T) {
	registry := annotation.NewRegistry(Get())

	assert.Empty(t, registry.ResolveAnnotations([]string{``}))
}
