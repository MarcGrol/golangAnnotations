package eventAnnotation

import "github.com/MarcGrol/golangAnnotations/annotation"

const (
	TypeEvent = "Event"

	ParamAggregate   = "aggregate"
	ParamIsRootEvent = "isrootevent"
)

// Register makes the annotation-registry aware of this annotation
func Register() {
	annotation.RegisterAnnotation(TypeEvent, []string{ParamAggregate, ParamIsRootEvent}, validateEventAnnotation)
}

func validateEventAnnotation(annot annotation.Annotation) bool {
	if annot.Name == TypeEvent {
		val, hasAggr := annot.Attributes[ParamAggregate]
		return hasAggr && val != ""
	}
	return false
}
