package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnnotationOk(t *testing.T) {
	s := Struct{
		DocLines:    []string{"Dummy", `    // {"Action":"Event","Data":{"Aggregate":"person"}}`},
		PackageName: "generator",
		Name:        "MyStruct",
	}
	assert.True(t, s.IsEvent())
	assert.Equal(t, "Person", s.GetAggregateName())
}

func TestAnnotationInvalid(t *testing.T) {
	s := Struct{
		DocLines:    []string{"// kvdkdakb"},
		PackageName: "generator",
		Name:        "MyStruct",
	}
	assert.False(t, s.IsEvent())
}
