package jsonAnnotation

import "github.com/MarcGrol/golangAnnotations/annotation"

const (
	TypeEnum      = "JsonEnum"
	TypeStruct    = "JsonStruct"
	ParamStripped = "stripped"
	ParamTolerant = "tolerant"
	ParamBase     = "base"
	ParamDefault  = "default"
)

func Get() []annotation.AnnotationDescriptor {
	return []annotation.AnnotationDescriptor{
		{
			Name:       TypeEnum,
			ParamNames: []string{ParamStripped, ParamTolerant, ParamBase, ParamDefault},
			Validator:  validateEnumAnnotation,
		},
		{
			Name:       TypeStruct,
			ParamNames: []string{},
			Validator:  validateStructAnnotation,
		}}
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
