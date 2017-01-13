package jsonAnnotation

import "github.com/MarcGrol/golangAnnotations/annotation"

const (
	typeEnum   = "JsonEnum"
	typeStruct = "JsonStruct"
)

// Register makes the annotation-registry aware of this annotation
func Register() {
	annotation.RegisterAnnotation(typeEnum, []string{}, validateEnumAnnotation)
	annotation.RegisterAnnotation(typeStruct, []string{}, validateStructAnnotation)
}

func validateEnumAnnotation(annot annotation.Annotation) bool {
	if annot.Name == typeEnum {
		return true
	}
	return false
}

func validateStructAnnotation(annot annotation.Annotation) bool {
	if annot.Name == typeStruct {
		return true
	}
	return false
}
