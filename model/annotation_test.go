package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGarbage(t *testing.T) {
	_, ok := resolveAnnotation(`// wvdwadbvb`)
	assert.False(t, ok)
}

func TestUnknownAction(t *testing.T) {
	_, ok := resolveAnnotation(`// {"Annotation":"Haha","With":{"X":"Y"}}`)
	assert.False(t, ok)
}

func TestCorrectEventAnnotation(t *testing.T) {
	annotation, ok := resolveAnnotation(`// {"Annotation":"Event","With":{"Aggregate":"Test"}}`)
	assert.True(t, ok)
	assert.Equal(t, "Event", annotation.Annotation)
	assert.Equal(t, "Test", annotation.With["Aggregate"])
}

func TestCorrectEventAnnotation2(t *testing.T) {
	aggregate, ok := resolveEventAnnotation([]string{`// {"Annotation":"Event","With":{"Aggregate":"test"}}`})
	assert.True(t, ok)
	assert.Equal(t, "Test", aggregate)
}

func TestIncompletegEventAnnotation(t *testing.T) {
	_, ok := resolveEventAnnotation([]string{`// {"Annotation":"Event"}`})
	assert.False(t, ok)
}

func TestCorrectRestServiceAnnotation(t *testing.T) {
	annotation, ok := resolveAnnotation(`// {"Annotation":"RestService","With":{"Path":"/person"}}`)
	assert.True(t, ok)
	assert.Equal(t, "RestService", annotation.Annotation)
	assert.Equal(t, "/person", annotation.With["Path"])
}

func TestCorrectRestServiceAnnotation2(t *testing.T) {
	path, ok := resolveRestServiceAnnotation([]string{`// {"Annotation":"RestService","With":{"Path":"/person"}}`})
	assert.True(t, ok)
	assert.Equal(t, "/person", path)
}

func TestIncompleteRestServiceAnnotation(t *testing.T) {
	_, ok := resolveAnnotation(`// {"Annotation":"RestService"}}`)
	assert.False(t, ok)
}

func TestCorrectRestOperationAnnotation(t *testing.T) {
	annotation, ok := resolveAnnotation(`// {"Annotation":"RestOperation","With":{"Method":"GET", "Path":"/person/:uid"}}`)
	assert.True(t, ok)
	assert.Equal(t, "RestOperation", annotation.Annotation)
	assert.Equal(t, "GET", annotation.With["Method"])
	assert.Equal(t, "/person/:uid", annotation.With["Path"])
}

func TestCorrectRestOperationAnnotation2(t *testing.T) {
	data, ok := resolveRestOperationAnnotation([]string{`// {"Annotation":"RestOperation","With":{"Method":"GET", "Path":"/person/:uid"}}`})
	assert.True(t, ok)
	assert.Equal(t, "GET", data["Method"])
	assert.Equal(t, "/person/:uid", data["Path"])
}

func TestIncompleteRestOperationAnnotation2(t *testing.T) {
	_, ok := resolveRestOperationAnnotation([]string{`// {"Annotation":"RestOperation","With":{"Method":"GET"}}`})
	assert.False(t, ok)

}
