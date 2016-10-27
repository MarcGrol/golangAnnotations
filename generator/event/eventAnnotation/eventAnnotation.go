package eventAnnotation

import "github.com/MarcGrol/golangAnnotations/annotation"

const (
	typeEvent      = "Event"
	paramAggregate = "aggregate"
	isRootEvent    = "isRootEvent"
)

// Register makes the annotation-registry aware of this annotation
func Register() {
	annotation.RegisterAnnotation(typeEvent, []string{paramAggregate, isRootEvent}, validateEventAnnotation)
}

func validateEventAnnotation(annot annotation.Annotation) bool {
	if annot.Name == typeEvent {
		val, hasAggr := annot.Attributes[paramAggregate]
		return (hasAggr && val != "")
	}
	return false
}
