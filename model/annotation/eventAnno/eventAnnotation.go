package eventAnno

import "github.com/MarcGrol/astTools/model/annotation"

const (
	typeEvent      = "Event"
	ParamAggregate = "aggregate"
)

// Register makes the annotation-registry aware of this annotation
func Register() {
	annotation.RegisterAnnotation(typeEvent, []string{ParamAggregate}, validateEventAnnotation)
}

func validateEventAnnotation(annot annotation.Annotation) bool {
	if annot.Name == typeEvent {
		val, hasAggr := annot.Attributes[ParamAggregate]
		return (hasAggr && val != "")
	}
	return false
}
