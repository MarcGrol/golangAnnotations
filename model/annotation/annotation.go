package annotation

import "strings"

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

type Annotation struct {
	Name       string
	Attributes map[string]string
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
		if err == nil {
			//log.Printf("*** Error parsing text-annotation %s: %+v", annotationDocline, err)
			// try the other format
			annotation, err = parseAnnotation(annotationDocline)
			if err != nil {
				//log.Printf("*** Error parsing json-annotation %s: %+v", annotationDocline, err)
				continue
			}
		}

		if annotation.Name != descriptor.name {
			//log.Printf("*** Annotation-line '%s' did NOT match %s", annotationDocline, descriptor.name)
			continue
		}

		ok := descriptor.validator(annotation)
		if !ok {
			//log.Printf("*** Annotation-line '%s' of type %s is invalid %+v", annotationDocline, descriptor.name, annotation)
			continue
		}
		//log.Printf("Valid %s annotation -line '%s'", annotation.Annotation, annotationDocline)

		return annotation, true
	}
	return Annotation{}, false
}
