package eventServiceAnnotation

import "github.com/MarcGrol/golangAnnotations/annotation"

const (
	typeEventService = "EventService"
	paramSelf = "self"
	paramAggregates = "aggregates"

	typeEventOperation = "EventOperation"
)

// Register makes the annotation-registry aware of these annotation
func Register() {
	annotation.RegisterAnnotation(typeEventService, []string{paramSelf,paramAggregates}, validateEventServiceAnnotation)
	annotation.RegisterAnnotation(typeEventOperation, []string{}, validateEventOperationAnnotation)
}

func validateEventServiceAnnotation(annot annotation.Annotation) bool {
	if annot.Name == typeEventService {
		_, ok := annot.Attributes[paramAggregates]
		return ok
	}
	return false
}

func validateEventOperationAnnotation(annot annotation.Annotation) bool {
	if annot.Name == typeEventOperation {
		return true
	}
	return false
}
