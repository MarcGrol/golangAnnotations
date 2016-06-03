package eventAnno

import (
	"github.com/MarcGrol/astTools/model/annotation"
)

const (
	TypeEvent      = "Event"
	ParamAggregate = "Aggregate"
)

// make annotation-registry aware of this annotation
func init() {
	annotation.RegisterAnnotation(TypeEvent, []string{ParamAggregate}, validateEventAnnotation)
}

func validateEventAnnotation(annot annotation.Annotation) bool {
	if annot.Annotation == TypeEvent {
		_, ok := annot.With[ParamAggregate]
		return ok
	}
	return false
}
