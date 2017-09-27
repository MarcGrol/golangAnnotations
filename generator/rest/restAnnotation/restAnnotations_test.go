package restAnnotation

import (
	"strings"
	"testing"

	"github.com/MarcGrol/golangAnnotations/annotation"
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

	assert.Empty(t, annotation.ResolveAnnotations([]string{`// @RestOperation()`}))
}

func TestPartialIncompleteRestOperationAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	assert.NotEmpty(t, annotation.ResolveAnnotations([]string{`// @RestOperation( Method = "GET")`}))
}

func TestPartialIncompleteRestOperationAnnotation2(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	assert.Empty(t, annotation.ResolveAnnotations([]string{`// @RestOperation( Path = "/foo")`}))
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

	ann, ok := annotation.ResolveAnnotationByName([]string{`// @RestOperation( Path="/foo", Method = "GET", optionalArgs="arg2,arg3", producesEvents="Order.OrderCreated, Basket.BasketFinalized")`}, "RestOperation")
	assert.True(t, ok)

	{
		optionalArgString, ok := ann.Attributes["optionalargs"]
		assert.True(t, ok)
		parts := strings.Split(optionalArgString, ",")
		assert.True(t, findArgInArray(parts, "arg2"))
		assert.True(t, findArgInArray(parts, "arg3"))
		assert.False(t, findArgInArray(parts, "arg1"))
	}

	{
		producesEventsArgString, ok := ann.Attributes["producesevents"]
		assert.True(t, ok)
		parts := strings.Split(producesEventsArgString, ",")
		assert.False(t, findArgInArray(parts, "vwfeweegw"))
		assert.True(t, findArgInArray(parts, "Order.OrderCreated"))
		assert.True(t, findArgInArray(parts, "Basket.BasketFinalized"))
	}
}

func TestCorrectRestServiceAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	ann, ok := annotation.ResolveAnnotationByName([]string{`// @RestService( Path = "/api")`}, "RestService")
	assert.True(t, ok)
	assert.Equal(t, "/api", ann.Attributes["path"])
}

func TestIncompleteRestServiceAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	assert.Empty(t, annotation.ResolveAnnotations([]string{`// @RestService()`}))
}

func TestEmptyRestServiceAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	assert.NotEmpty(t, annotation.ResolveAnnotations([]string{`// @RestService( Path = "")`}))
}
