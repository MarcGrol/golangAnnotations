package annotation

import "strings"

type Annotation struct {
	Name       string
	Attributes map[string]string
}

type ValidationFunc func(annot Annotation) bool

type annotationDescriptor struct {
	name       string
	paramNames []string
	validator  ValidationFunc
}

var annotationRegistry []annotationDescriptor = []annotationDescriptor{}

func ClearRegisteredAnnotations() {
	annotationRegistry = []annotationDescriptor{}
}

func RegisterAnnotation(name string, paramNames []string, validator ValidationFunc) {
	annotationRegistry = append(annotationRegistry, annotationDescriptor{name: name, paramNames: paramNames, validator: validator})
}

func ResolveAnnotations(annotationDocline []string) (Annotation, bool) {
	for _, line := range annotationDocline {
		a, ok := ResolveAnnotation(strings.TrimSpace(line))
		if ok {
			return a, ok
		}
	}
	return Annotation{}, false
}

func ResolveAnnotation(annotationDocline string) (Annotation, bool) {
	for _, descriptor := range annotationRegistry {
		annotation, err := parseAnnotation(annotationDocline)
		annotation, err = parseAnnotation(annotationDocline)
		if err != nil {
			continue
		}

		if annotation.Name != descriptor.name {
			continue
		}

		ok := descriptor.validator(annotation)
		if !ok {
			continue
		}

		return annotation, true
	}
	return Annotation{}, false
}
