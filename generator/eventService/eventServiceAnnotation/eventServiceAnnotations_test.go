package eventServiceAnnotation

import (
	"strings"
	"testing"

	"github.com/MarcGrol/golangAnnotations/annotation"
	"github.com/stretchr/testify/assert"
)

func TestIncompleteRestServiceAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	assert.Empty(t, annotation.ResolveAnnotations([]string{`// @EventService()`}))
}

func TestCorrectEventServiceAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	ann, ok := annotation.ResolveAnnotationByName([]string{`// @EventService( Self = "caregiverService", Subscriptions = "order,basket")`}, "EventService")
	assert.True(t, ok)

	self, ok := ann.Attributes["self"]
	assert.True(t, ok)
	assert.Equal(t, "caregiverService", self)

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
