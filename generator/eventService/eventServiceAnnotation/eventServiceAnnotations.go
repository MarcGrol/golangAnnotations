package eventServiceAnnotation

import "github.com/MarcGrol/golangAnnotations/annotation"

const (
	typeEventService = "EventService"
	paramSelf = "self"
	paramSubscriptions = "subscriptions"

	typeEventOperation = "EventOperation"
)

// Register makes the annotation-registry aware of these annotation
func Register() {
	annotation.RegisterAnnotation(typeEventService, []string{paramSelf, paramSubscriptions}, validateEventServiceAnnotation)
	annotation.RegisterAnnotation(typeEventOperation, []string{}, validateEventOperationAnnotation)
}

func validateEventServiceAnnotation(annot annotation.Annotation) bool {
	if annot.Name == typeEventService {
		_, ok := annot.Attributes[paramSubscriptions]
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
