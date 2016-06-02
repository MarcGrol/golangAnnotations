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

type annotation struct {
	name       string
	paramNames []string
}

var annotations []annotation = []annotation{
	{name: annotationTypeEvent, paramNames: []string{"Aggregate"}},
	{name: annotationTypeRestService, paramNames: []string{"Path"}},
	{name: annotationTypeRestOperation, paramNames: []string{"Method", "Path"}},
}

func resolveEventAnnotation(lines []string) (string, bool) {
	for _, line := range lines {
		annotation, ok := resolveAnnotation(strings.TrimSpace(line))
		if ok && annotation.Action == annotationTypeEvent {
			_, ok := annotation.Data["Aggregate"]
			return toFirstUpper(annotation.Data["Aggregate"]), ok
		}
	}
	return "", false
}

func resolveRestServiceAnnotation(lines []string) (string, bool) {
	for _, line := range lines {
		annotation, ok := resolveAnnotation(strings.TrimSpace(line))
		if ok && annotation.Action == annotationTypeRestService {
			_, ok := annotation.Data["Path"]
			return annotation.Data["Path"], ok
		}
	}
	return "", false
}

func resolveRestOperationAnnotation(lines []string) (map[string]string, bool) {
	for _, line := range lines {
		annotation, ok := resolveAnnotation(strings.TrimSpace(line))
		if ok && annotation.Action == annotationTypeRestOperation {
			_, hasPath := annotation.Data["Path"]
			_, hasMethod := annotation.Data["Method"]
			return annotation.Data, hasPath && hasMethod
		}
	}
	return map[string]string{}, false
}

type Annotation struct {
	Action string
	Data   map[string]string
}

func resolveAnnotation(annotationDocline string) (Annotation, bool) {
	for _, ann := range annotations {
		annotation, err := parseAnnotation(annotationDocline)
		if err != nil {
			log.Printf("*** Error unmarshalling RestOperationAnnotation %s: %+v", annotationDocline, err)
			continue
		}
		if annotation.Action == ann.name {
			log.Printf("Annotation-line '%s' MATCHED %+v -> %+v", annotationDocline, ann, annotation)
			return annotation, true
		} else {
			log.Printf("Annotation-line '%s' did NOT match %+v", annotationDocline, ann)
		}
	}
	return Annotation{}, false
}

func parseAnnotation(annotationDocline string) (Annotation, error) {
	withoutComment := strings.TrimLeft(strings.TrimSpace(annotationDocline), "/")

	var annotation Annotation
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
