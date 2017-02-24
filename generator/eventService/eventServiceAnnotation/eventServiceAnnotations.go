package eventServiceAnnotation

import "github.com/MarcGrol/golangAnnotations/annotation"

const (
	TypeEventService   = "EventService"
	TypeEventOperation = "EventOperation"
	ParamSelf          = "self"
	ParamSubscriptions = "subscriptions"
	ParamAsync         = "async"
)

// Register makes the annotation-registry aware of these annotation
func Register() {
	annotation.RegisterAnnotation(TypeEventService, []string{ParamSelf, ParamSubscriptions, ParamAsync}, validateEventServiceAnnotation)
	annotation.RegisterAnnotation(TypeEventOperation, []string{}, validateEventOperationAnnotation)
}

func validateEventServiceAnnotation(annot annotation.Annotation) bool {
	if annot.Name == TypeEventService {
		_, ok := annot.Attributes[ParamSubscriptions]
		return ok
	}
	return false
}

func validateEventOperationAnnotation(annot annotation.Annotation) bool {
	if annot.Name == TypeEventOperation {
		return true
	}
	return false
}
