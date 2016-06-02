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
	_, ok := resolveAnnotation(`// {"Action":"Haha","Data":{"X":"Y"}}`)
	assert.False(t, ok)
}

func TestCorrectEventAnnotation(t *testing.T) {
	annotation, ok := resolveAnnotation(`// {"Action":"Event","Data":{"Aggregate":"Test"}}`)
	assert.True(t, ok)
	assert.Equal(t, "Event", annotation.Action)
	assert.Equal(t, "Test", annotation.Data["Aggregate"])
}

func TestCorrectEventAnnotation2(t *testing.T) {
	aggregate, ok := resolveEventAnnotation([]string{`// {"Action":"Event","Data":{"Aggregate":"test"}}`})
	assert.True(t, ok)
	assert.Equal(t, "Test", aggregate)
}

func TestIncompletegEventAnnotation(t *testing.T) {
	_, ok := resolveEventAnnotation([]string{`// {"Action":"Event"}`})
	assert.False(t, ok)
}

func TestCorrectRestServiceAnnotation(t *testing.T) {
	annotation, ok := resolveAnnotation(`// {"Action":"RestService","Data":{"Path":"/person"}}`)
	assert.True(t, ok)
	assert.Equal(t, "RestService", annotation.Action)
	assert.Equal(t, "/person", annotation.Data["Path"])
}

func TestCorrectRestServiceAnnotation2(t *testing.T) {
	path, ok := resolveRestServiceAnnotation([]string{`// {"Action":"RestService","Data":{"Path":"/person"}}`})
	assert.True(t, ok)
	assert.Equal(t, "/person", path)
}

func TestIncompleteRestServiceAnnotation(t *testing.T) {
	_, ok := resolveAnnotation(`// {"Action":"RestService"}}`)
	assert.False(t, ok)
}

func TestCorrectRestOperationAnnotation(t *testing.T) {
	annotation, ok := resolveAnnotation(`// {"Action":"RestOperation","Data":{"Method":"GET", "Path":"/person/:uid"}}`)
	assert.True(t, ok)
	assert.Equal(t, "RestOperation", annotation.Action)
	assert.Equal(t, "GET", annotation.Data["Method"])
	assert.Equal(t, "/person/:uid", annotation.Data["Path"])
}

func TestCorrectRestOperationAnnotation2(t *testing.T) {
	data, ok := resolveRestOperationAnnotation([]string{`// {"Action":"RestOperation","Data":{"Method":"GET", "Path":"/person/:uid"}}`})
	assert.True(t, ok)
	assert.Equal(t, "GET", data["Method"])
	assert.Equal(t, "/person/:uid", data["Path"])
}

func TestIncompleteRestOperationAnnotation2(t *testing.T) {
	_, ok := resolveRestOperationAnnotation([]string{`// {"Action":"RestOperation","Data":{"Method":"GET"}}`})
	assert.False(t, ok)

}
