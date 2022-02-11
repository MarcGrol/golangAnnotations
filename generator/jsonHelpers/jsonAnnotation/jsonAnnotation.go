package jsonAnnotation

import "github.com/f0rt/golangAnnotations/generator/annotation"

const (
	TypeEnum      = "JsonEnum"
	TypeStruct    = "JsonStruct"
	ParamStripped = "stripped"
	ParamLiteral  = "literal"
	ParamTolerant = "tolerant"
	ParamBase     = "base"
	ParamDefault  = "default"
)

func Get() []annotation.AnnotationDescriptor {
	return []annotation.AnnotationDescriptor{
		{
			Name:       TypeEnum,
			ParamNames: []string{ParamStripped, ParamLiteral, ParamTolerant, ParamBase, ParamDefault},
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
