package restAnnotation

import (
	"strings"
	"testing"

	"github.com/f0rt/golangAnnotations/generator/annotation"
	"github.com/stretchr/testify/assert"
)

func TestCorrectRestOperationAnnotation(t *testing.T) {
	registry := annotation.NewRegistry(Get())

	a, ok := registry.ResolveAnnotation(`// @RestOperation( Method = "GET", path = "/person/:uid" )`)
	assert.True(t, ok)
	assert.Equal(t, "GET", a.Attributes["method"])
	assert.Equal(t, "/person/:uid", a.Attributes["path"])
}

func TestIncompleteRestOperationAnnotation(t *testing.T) {
	registry := annotation.NewRegistry(Get())

	assert.Empty(t, registry.ResolveAnnotations([]string{`// @RestOperation()`}))
}

func TestPartialIncompleteRestOperationAnnotation(t *testing.T) {
	registry := annotation.NewRegistry(Get())

	assert.NotEmpty(t, registry.ResolveAnnotations([]string{`// @RestOperation( Method = "GET")`}))
}

func TestPartialIncompleteRestOperationAnnotation2(t *testing.T) {
	registry := annotation.NewRegistry(Get())

	assert.Empty(t, registry.ResolveAnnotations([]string{`// @RestOperation( Path = "/foo")`}))
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
	registry := annotation.NewRegistry(Get())

	ann, ok := registry.ResolveAnnotationByName([]string{`// @RestOperation( Path="/foo", Method = "GET", optionalArgs="arg2,arg3", producesEvents="Order.OrderCreated, Basket.BasketFinalized")`}, "RestOperation")
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
	registry := annotation.NewRegistry(Get())

	ann, ok := registry.ResolveAnnotationByName([]string{`// @RestService( Path = "/api")`}, "RestService")
	assert.True(t, ok)
	assert.Equal(t, "/api", ann.Attributes["path"])
}

func TestIncompleteRestServiceAnnotation(t *testing.T) {
	registry := annotation.NewRegistry(Get())

	assert.Empty(t, registry.ResolveAnnotations([]string{`// @RestService()`}))
}

func TestEmptyRestServiceAnnotation(t *testing.T) {
	registry := annotation.NewRegistry(Get())

	assert.NotEmpty(t, registry.ResolveAnnotations([]string{`// @RestService( Path = "")`}))
}
