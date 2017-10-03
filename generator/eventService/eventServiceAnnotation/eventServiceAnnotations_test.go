package eventServiceAnnotation

import (
	"testing"

	"github.com/MarcGrol/golangAnnotations/annotation"
	"github.com/stretchr/testify/assert"
)

func TestCorrectEventServiceAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	ann, ok := annotation.ResolveAnnotationByName([]string{`// @EventService( Self = "caregiverService", async="true", admin="true", producesEvents="x,y" )`}, "EventService")
	assert.True(t, ok)
	{
		self, ok := ann.Attributes[ParamSelf]
		assert.True(t, ok)
		assert.Equal(t, "caregiverService", self)
	}
	{
		async, ok := ann.Attributes[ParamAsync]
		assert.True(t, ok)
		assert.Equal(t, "true", async)
	}
	{
		admin, ok := ann.Attributes[ParamAdmin]
		assert.True(t, ok)
		assert.Equal(t, "true", admin)
	}
	{
		producesEvents, ok := ann.Attributes[ParamProducesEvents]
		assert.True(t, ok)
		assert.Equal(t, "x,y", producesEvents)
	}
}

func TestCorrectEventOperationAnnotation(t *testing.T) {
	annotation.ClearRegisteredAnnotations()
	Register()

	_, ok := annotation.ResolveAnnotation(`// @EventOperation( topic = "order" )`)
	assert.True(t, ok)
}
