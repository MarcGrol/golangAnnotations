package eventAnnotation

import "github.com/MarcGrol/golangAnnotations/annotation"

const (
	TypeEvent        = "Event"
	ParamAggregate   = "aggregate"
	ParamIsRootEvent = "isrootevent"
	ParamIsTransient = "istransient"
)

// Register makes the annotation-registry aware of this annotation
func Get() []annotation.AnnotationDescriptor {
	return []annotation.AnnotationDescriptor{
		{
			Name:       TypeEvent,
			ParamNames: []string{ParamAggregate, ParamIsRootEvent, ParamIsTransient},
			Validator:  validateEventAnnotation,
		},
	}
}

func validateEventAnnotation(annot annotation.Annotation) bool {
	if annot.Name == TypeEvent {
		val, hasAggr := annot.Attributes[ParamAggregate]
		return hasAggr && val != ""
	}
	return false
}
