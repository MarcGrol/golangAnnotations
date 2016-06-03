package eventAnno

import (
	"testing"

	"github.com/MarcGrol/astTools/model"
	"github.com/MarcGrol/astTools/model/annotation"
	"github.com/stretchr/testify/assert"
)

func TestAnnotationInvalid(t *testing.T) {
	s := model.Struct{
		DocLines:    []string{"// kvdkdakb"},
		PackageName: "generator",
		Name:        "MyStruct",
	}
	assert.False(t, s.IsEvent())
}

func TestAnnotationOk(t *testing.T) {
	s := model.Struct{
		DocLines:    []string{"Dummy", `    // {"Annotation":"Event","With":{"Aggregate":"person"}}`},
		PackageName: "generator",
		Name:        "MyStruct",
	}
	assert.True(t, s.IsEvent())
	assert.Equal(t, "person", s.GetAggregateName())
}

func TestCorrectEventAnnotation(t *testing.T) {
	annot, ok := annotation.ResolveAnnotations([]string{`// {"Annotation":"Event","With":{"Aggregate":"test"}}`})
	assert.True(t, ok)
	assert.Equal(t, "test", annot.With["Aggregate"])
}

func TestIncompleteEventAnnotation(t *testing.T) {
	_, ok := annotation.ResolveAnnotations([]string{`// {"Annotation":"Event"}`})
	assert.False(t, ok)
}
