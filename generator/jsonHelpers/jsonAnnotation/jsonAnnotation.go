package jsonAnnotation

import "github.com/MarcGrol/golangAnnotations/annotation"

const (
	TypeEnum   = "JsonEnum"
	TypeStruct = "JsonStruct"
	ParamBase  = "base"
)

// Register makes the annotation-registry aware of this annotation
func Register() {
	annotation.RegisterAnnotation(TypeEnum, []string{ParamBase}, validateEnumAnnotation)
	annotation.RegisterAnnotation(TypeStruct, []string{}, validateStructAnnotation)
}

func validateEnumAnnotation(annot annotation.Annotation) bool {
	if annot.Name == TypeEnum {
		return true
	}
	return false
}

func validateStructAnnotation(annot annotation.Annotation) bool {
	if annot.Name == TypeStruct {
		return true
	}
	return false
}
