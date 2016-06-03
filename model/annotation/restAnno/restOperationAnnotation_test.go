package restAnno

import (
	"testing"

	"github.com/MarcGrol/astTools/model/annotation"
	"github.com/stretchr/testify/assert"
)

func TestCorrectRestOperationAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	a, ok := annotation.ResolveAnnotation(`// {"Annotation":"RestOperation","With":{"Method":"GET", "Path":"/person/:uid"}}`)
	assert.True(t, ok)
	assert.Equal(t, "GET", a.With["Method"])
	assert.Equal(t, "/person/:uid", a.With["Path"])
}

func TestIncompleteRestOperationAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	_, ok := annotation.ResolveAnnotations([]string{`// {"Annotation":"RestOperation","With":{"Method":"GET"}}`})
	assert.False(t, ok)
}

func TestCorrectRestServiceAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	a, ok := annotation.ResolveAnnotations([]string{`// {"Annotation":"RestService","With":{"Path":"/person"}}`})
	assert.True(t, ok)
	assert.Equal(t, "/person", a.With["Path"])
}

func TestIncompleteRestServiceAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	_, ok := annotation.ResolveAnnotations([]string{`// {"Annotation":"RestService"`})
	assert.False(t, ok)
}
