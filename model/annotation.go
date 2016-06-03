package model

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

const (
	annotationTypeEvent         = "Event"
	annotationTypeRestService   = "RestService"
	annotationTypeRestOperation = "RestOperation"
)

type annotationDescriptor struct {
	name       string
	paramNames []string
}

var annotations []annotationDescriptor = []annotationDescriptor{
	{name: annotationTypeEvent, paramNames: []string{"Aggregate"}},
	{name: annotationTypeRestService, paramNames: []string{"Path"}},
	{name: annotationTypeRestOperation, paramNames: []string{"Method", "Path"}},
}

func resolveEventAnnotation(lines []string) (string, bool) {
	for _, line := range lines {
		annotation, ok := resolveAnnotation(strings.TrimSpace(line))
		if ok && annotation.Annotation == annotationTypeEvent {
			_, ok := annotation.With["Aggregate"]
			return toFirstUpper(annotation.With["Aggregate"]), ok
		}
	}
	return "", false
}

func resolveRestServiceAnnotation(lines []string) (string, bool) {
	for _, line := range lines {
		annotation, ok := resolveAnnotation(strings.TrimSpace(line))
		if ok && annotation.Annotation == annotationTypeRestService {
			_, ok := annotation.With["Path"]
			return annotation.With["Path"], ok
		}
	}
	return "", false
}

func resolveRestOperationAnnotation(lines []string) (map[string]string, bool) {
	for _, line := range lines {
		annotation, ok := resolveAnnotation(strings.TrimSpace(line))
		if ok && annotation.Annotation == annotationTypeRestOperation {
			_, hasPath := annotation.With["Path"]
			_, hasMethod := annotation.With["Method"]
			return annotation.With, hasPath && hasMethod
		}
	}
	return map[string]string{}, false
}

type AnnotationLine struct {
	Annotation string
	With       map[string]string
}

func resolveAnnotation(annotationDocline string) (AnnotationLine, bool) {
	for _, ann := range annotations {
		annotation, err := parseAnnotation(annotationDocline)
		if err != nil {
			log.Printf("*** Error unmarshalling RestOperationAnnotation %s: %+v", annotationDocline, err)
			continue
		}
		if annotation.Annotation == ann.name {
			log.Printf("Annotation-line '%s' MATCHED %+v -> %+v", annotationDocline, ann, annotation)
			return annotation, true
		} else {
			log.Printf("Annotation-line '%s' did NOT match %+v", annotationDocline, ann)
		}
	}
	return AnnotationLine{}, false
}

func parseAnnotation(annotationDocline string) (AnnotationLine, error) {
	withoutComment := strings.TrimLeft(strings.TrimSpace(annotationDocline), "/")

	var annotation AnnotationLine
	err := json.Unmarshal([]byte(withoutComment), &annotation)
	if err != nil {
		return annotation, err
	}
	return annotation, nil
}

func toFirstUpper(in string) string {
	if len(in) == 0 {
		return in
	}
	return strings.ToUpper(fmt.Sprintf("%c", in[0])) + in[1:]
}
