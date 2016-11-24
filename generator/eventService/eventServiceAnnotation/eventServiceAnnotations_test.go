package eventServiceAnnotation

import (
	"github.com/MarcGrol/golangAnnotations/annotation"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestIncompleteRestServiceAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	_, ok := annotation.ResolveAnnotations([]string{`// @EventService()`})
	assert.False(t, ok)
}

func TestCorrectEventServiceAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	ann, ok := annotation.ResolveAnnotations([]string{`// @EventService( Self = "caregiver", Subscriptions = "order,basket")`})
	assert.True(t, ok)

	self, ok := ann.Attributes["self"]
	assert.True(t, ok)
	assert.Equal(t, "caregiver", self)

	aggregates, ok := ann.Attributes["subscriptions"]
	assert.True(t, ok)
	parts := strings.Split(aggregates, ",")
	assert.True(t, findArgInArray(parts, "order"))
	assert.True(t, findArgInArray(parts, "basket"))
	assert.False(t, findArgInArray(parts, "caregiver"))
}

func findArgInArray(array []string, toMatch string) bool {
	for _, p := range array {
		if strings.Trim(p, " ") == toMatch {
			return true
		}
	}
	return false
}

func TestCorrectEventOperationAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	_, ok := annotation.ResolveAnnotation(`// @EventOperation( )`)
	assert.True(t, ok)
}
