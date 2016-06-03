package restAnno

import (
	"testing"

	"github.com/MarcGrol/astTools/model/annotation"
	"github.com/stretchr/testify/assert"
)

func TestCorrectRestServiceAnnotation(t *testing.T) {
	a, ok := annotation.ResolveAnnotations([]string{`// {"Annotation":"RestService","With":{"Path":"/person"}}`})
	assert.True(t, ok)
	assert.Equal(t, "/person", a.With["Path"])
}

func TestIncompleteRestServiceAnnotation(t *testing.T) {
	_, ok := annotation.ResolveAnnotations([]string{`// {"Annotation":"RestService"`})
	assert.False(t, ok)
}
