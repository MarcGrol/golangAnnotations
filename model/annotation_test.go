package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEventAnnotation(t *testing.T) {
	annoType, m := resolveAnnotation("// +event -> aggregate: test")
	assert.Equal(t, annoType, annotationTypeEvent)
	assert.Equal(t, "test", m["aggregate"])
}

func TestWrongAnnotation(t *testing.T) {
	annoType, _ := resolveAnnotation("// event -> aggregate: Test")
	assert.Equal(t, annoType, annotationTypeUnknown)
}
