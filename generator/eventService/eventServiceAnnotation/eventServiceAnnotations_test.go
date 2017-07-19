package eventServiceAnnotation

import (
	"testing"

	"github.com/MarcGrol/golangAnnotations/annotation"
	"github.com/stretchr/testify/assert"
)

func TestCorrectEventServiceAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	ann, ok := annotation.ResolveAnnotationByName([]string{`// @EventService( Self = "caregiverService" )`}, "EventService")
	assert.True(t, ok)

	self, ok := ann.Attributes["self"]
	assert.True(t, ok)
	assert.Equal(t, "caregiverService", self)
}

func TestCorrectEventOperationAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	_, ok := annotation.ResolveAnnotation(`// @EventOperation( topic = "order" )`)
	assert.True(t, ok)
}
