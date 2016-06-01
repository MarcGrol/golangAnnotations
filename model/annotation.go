package model

import (
	"fmt"
	"log"
	"strings"
)

type annotationType int

const (
	annotationTypeUnknown annotationType = iota
	annotationTypeEvent   annotationType = iota + 1
)

type ParseFunc func(annotation, string) (map[string]string, bool)

type annotation struct {
	name       string
	annoType   annotationType
	format     string
	paramNames []string
	parseFunc  ParseFunc
}

var annotations []annotation = []annotation{
	{name: "event", annoType: annotationTypeEvent, format: "// +event -> aggregate: %s", paramNames: []string{"aggregate"}, parseFunc: parseEventAnnotation},
}

func resolveEventAnnotation(lines []string) (string, bool) {
	for _, line := range lines {
		t, m := resolveAnnotation(strings.TrimSpace(line))
		if t == annotationTypeEvent {
			val, ok := m["aggregate"]
			return toFirstUpper(val), ok
		}
	}
	return "", false
}

func parseEventAnnotation(ann annotation, annotationDocline string) (map[string]string, bool) {
	matched := false
	annotationData := make(map[string]string)
	aggregateName := ""

	count, err := fmt.Sscanf(annotationDocline, ann.format, &aggregateName)
	if err == nil && count == len(ann.paramNames) {
		matched = true
		annotationData["aggregate"] = aggregateName
	}
	return annotationData, matched
}

func resolveAnnotation(annotationDocline string) (annotationType, map[string]string) {
	annotationType := annotationTypeUnknown
	annotationData := make(map[string]string)

	for _, ann := range annotations {
		var ok bool
		annotationData, ok = ann.parseFunc(ann, annotationDocline)
		if ok {
			annotationType = ann.annoType
			log.Printf("Annotation-line '%s' matched %+v", annotationDocline, ann)
			break
		} else {
			log.Printf("Annotation-line '%s' did not match %+v", annotationDocline, ann)
		}
	}
	return annotationType, annotationData
}

func toFirstUpper(in string) string {
	return strings.ToUpper(fmt.Sprintf("%c", in[0])) + in[1:]
}
