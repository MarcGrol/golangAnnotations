package annotation

import "strings"

type AnnotationRegister interface {
	ResolveAnnotations(annotationDocline []string) []Annotation
	ResolveAnnotationByName(annotationDocline []string, name string) (Annotation, bool)
	ResolveAnnotation(annotationDocline string) (Annotation, bool)
}

type annotationRegistry struct {
	descriptors []AnnotationDescriptor
}

func NewRegistry(descriptors []AnnotationDescriptor) AnnotationRegister {
	return &annotationRegistry{
		descriptors: descriptors,
	}
}

type Annotation struct {
	Name       string
	Attributes map[string]string
}

type validationFunc func(annot Annotation) bool

type AnnotationDescriptor struct {
	Name       string
	ParamNames []string
	Validator  validationFunc
}

func (ar *annotationRegistry) ResolveAnnotations(annotationDocline []string) []Annotation {
	annotations := make([]Annotation, 0)
	for _, line := range annotationDocline {
		if ann, ok := ar.ResolveAnnotation(strings.TrimSpace(line)); ok {
			annotations = append(annotations, ann)
		}
	}
	return annotations
}

func (ar *annotationRegistry) ResolveAnnotationByName(annotationDocline []string, name string) (Annotation, bool) {
	for _, line := range annotationDocline {
		ann, ok := ar.ResolveAnnotation(strings.TrimSpace(line))
		if ok && ann.Name == name {
			return ann, true
		}
	}
	return Annotation{}, false
}

func (ar *annotationRegistry) ResolveAnnotation(annotationDocline string) (Annotation, bool) {
	for _, descriptor := range ar.descriptors {
		ann, err := parseAnnotation(annotationDocline)
		if err != nil {
			continue
		}

		if ann.Name != descriptor.Name {
			continue
		}

		ok := descriptor.Validator(ann)
		if !ok {
			continue
		}

		return ann, true
	}
	return Annotation{}, false
}
