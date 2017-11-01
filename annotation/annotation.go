package annotation

import "strings"

type Annotation struct {
	Name       string
	Attributes map[string]string
}

type validationFunc func(annot Annotation) bool

type annotationDescriptor struct {
	name       string
	paramNames []string
	validator  validationFunc
}

var annotationRegistry []annotationDescriptor = []annotationDescriptor{}

func ClearRegisteredAnnotations() {
	annotationRegistry = []annotationDescriptor{}
}

func RegisterAnnotation(name string, paramNames []string, validator validationFunc) {
	annotationRegistry = append(annotationRegistry, annotationDescriptor{name: name, paramNames: paramNames, validator: validator})
}

func ResolveAnnotations(annotationDocline []string) []Annotation {
	annotations := []Annotation{}
	for _, line := range annotationDocline {
		ann, ok := ResolveAnnotation(strings.TrimSpace(line))
		if ok {
			annotations = append(annotations, ann)
		}
	}
	return annotations
}

func ResolveAnnotationByName(annotationDocline []string, name string) (Annotation, bool) {
	for _, line := range annotationDocline {
		ann, ok := ResolveAnnotation(strings.TrimSpace(line))
		if ok && ann.Name == name {
			return ann, true
		}
	}
	return Annotation{}, false
}

func ResolveAnnotation(annotationDocline string) (Annotation, bool) {
	for _, descriptor := range annotationRegistry {
		ann, err := parseAnnotation(annotationDocline)
		if err != nil {
			continue
		}

		if ann.Name != descriptor.name {
			continue
		}

		ok := descriptor.validator(ann)
		if !ok {
			continue
		}

		return ann, true
	}
	return Annotation{}, false
}
