package eventAnnotation

import "github.com/MarcGrol/golangAnnotations/generator/annotation"

const (
	TypeEvent         = "Event"
	TypeEventPart     = "EventPart"
	ParamAggregate    = "aggregate"
	ParamIsRootEvent  = "isrootevent"
	ParamIsTransient  = "istransient"
	ParamIsSensitive  = "issensitive"
	FieldTagSensitive = "sensitive"
)

// Register makes the annotation-registry aware of this annotation
func Get() []annotation.AnnotationDescriptor {
	return []annotation.AnnotationDescriptor{
		{
			Name:       TypeEvent,
			ParamNames: []string{ParamAggregate, ParamIsRootEvent, ParamIsTransient, ParamIsSensitive},
			Validator:  validateEventAnnotation,
		},
		{
			Name:       TypeEventPart,
			ParamNames: []string{ParamIsSensitive},
			Validator:  validateEventAnnotation,
		},
	}
}

func validateEventAnnotation(annot annotation.Annotation) bool {
	switch annot.Name {
	case TypeEvent:
		val, hasAggr := annot.Attributes[ParamAggregate]
		return hasAggr && val != ""
	case TypeEventPart:
		return true
	}
	return false
}
