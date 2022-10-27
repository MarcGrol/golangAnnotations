package eventServiceAnnotation

import "github.com/f0rt/golangAnnotations/generator/annotation"

const (
	TypeEventService    = "EventService"
	TypeEventOperation  = "EventOperation"
	ParamSelf           = "self"
	ParamTopic          = "topic"
	ParamProcess        = "process"
	ParamDelayed        = "delayed"
	ParamNoTest         = "notest"
	ParamProducesEvents = "producesevents"
)

func Get() []annotation.AnnotationDescriptor {
	return []annotation.AnnotationDescriptor{
		{
			Name:       TypeEventService,
			ParamNames: []string{ParamSelf, ParamNoTest},
			Validator:  validateEventServiceAnnotation,
		},
		{
			Name:       TypeEventOperation,
			ParamNames: []string{ParamTopic, ParamProcess, ParamDelayed},
			Validator:  validateEventOperationAnnotation,
		}}
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
