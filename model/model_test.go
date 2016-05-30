package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnnotationOk(t *testing.T) {
	s := Struct{
		DocLines:    []string{"Dummy", "   // +event -> aggregate: person   "},
		PackageName: "generator",
		Name:        "MyStruct",
	}
	assert.True(t, s.IsEvent())
	assert.Equal(t, "person", s.GetAggregateName())
}

func TestAnnotationInvalid(t *testing.T) {
	s := Struct{
		DocLines:    []string{"// kvdkdakb"},
		PackageName: "generator",
		Name:        "MyStruct",
	}
	assert.False(t, s.IsEvent())
}
