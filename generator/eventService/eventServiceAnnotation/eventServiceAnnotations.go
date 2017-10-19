package eventServiceAnnotation

import "github.com/MarcGrol/golangAnnotations/annotation"

const (
	TypeEventService    = "EventService"
	TypeEventOperation  = "EventOperation"
	ParamSelf           = "self"
	ParamAsync          = "async"
	ParamTopic          = "topic"
	ParamProcess        = "process"
	ParamNoTest         = "notest"
	ParamProducesEvents = "producesevents"
)

// Register makes the annotation-registry aware of these annotation
func Register() {
	annotation.RegisterAnnotation(TypeEventService, []string{ParamSelf, ParamAsync, ParamNoTest}, validateEventServiceAnnotation)
	annotation.RegisterAnnotation(TypeEventOperation, []string{ParamTopic, ParamProcess}, validateEventOperationAnnotation)
}

func validateEventServiceAnnotation(annot annotation.Annotation) bool {
	if annot.Name == TypeEventService {
		return true
	}
	return false
}

func validateEventOperationAnnotation(annot annotation.Annotation) bool {
	if annot.Name == TypeEventOperation {
		_, ok := annot.Attributes[ParamTopic]
		return ok
	}
	return false
}
