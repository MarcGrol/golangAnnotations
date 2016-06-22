package restAnno

import (
	"testing"

	"github.com/MarcGrol/astTools/model/annotation"
	"github.com/stretchr/testify/assert"
)

func TestCorrectRestOperationAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	a, ok := annotation.ResolveAnnotation(`// @RestOperation( Method = "GET", path = "/person/:uid" )`)
	assert.True(t, ok)
	assert.Equal(t, "GET", a.Attributes["method"])
	assert.Equal(t, "/person/:uid", a.Attributes["path"])
}

func TestIncompleteRestOperationAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	_, ok := annotation.ResolveAnnotations([]string{`// @RestOperation()`})
	assert.False(t, ok)
}

func TestPartialIncompleteRestOperationAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	_, ok := annotation.ResolveAnnotations([]string{`// @RestOperation( Method = "GET")`})
	assert.False(t, ok)
}

func TestPartialIncompleteRestOperationAnnotation2(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	_, ok := annotation.ResolveAnnotations([]string{`// @RestOperation( Path = "/foo")`})
	assert.False(t, ok)
}

func TestCorrectRestServiceAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	a, ok := annotation.ResolveAnnotations([]string{`// @RestService( Path = "/api")`})
	assert.True(t, ok)
	assert.Equal(t, "/api", a.Attributes["path"])
}

func TestIncompleteRestServiceAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	_, ok := annotation.ResolveAnnotations([]string{`// @RestService()`})
	assert.False(t, ok)
}

func TestEmptyRestServiceAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	_, ok := annotation.ResolveAnnotations([]string{`// @RestService( Path = "")`})
	assert.True(t, ok)
}
