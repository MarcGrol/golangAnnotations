package enumAnnotation

import "github.com/MarcGrol/golangAnnotations/annotation"

const (
	typeEnum = "Enum"
)

// Register makes the annotation-registry aware of this annotation
func Register() {
	annotation.RegisterAnnotation(typeEnum, []string{}, validateEnumAnnotation)
}

func validateEnumAnnotation(annot annotation.Annotation) bool {
	if annot.Name == typeEnum {
		return true
	}
	return false
}
