package restAnnotation

import (
	"github.com/MarcGrol/golangAnnotations/annotation"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
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
	assert.True(t, ok)
}

func TestPartialIncompleteRestOperationAnnotation2(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	_, ok := annotation.ResolveAnnotations([]string{`// @RestOperation( Path = "/foo")`})
	assert.False(t, ok)
}

func findArgInArray(array []string, toMatch string) bool {
	for _, p := range array {
		if strings.Trim(p, " ") == toMatch {
			return true
		}
	}
	return false
}

func TestOptionalArgRestOperationAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	ann, ok := annotation.ResolveAnnotations([]string{`// @RestOperation( Path="/foo", Method = "GET", optionalArgs="arg2,arg3")`})

	assert.True(t, ok)
	optionalArgString, ok := ann.Attributes["optionalargs"]
	assert.True(t, ok)
	parts := strings.Split(optionalArgString, ",")
	assert.True(t, findArgInArray(parts, "arg2"))
	assert.True(t, findArgInArray(parts, "arg3"))
	assert.False(t, findArgInArray(parts, "arg1"))
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
