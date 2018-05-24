package annotation

import (
	"fmt"
	"strings"
	"text/scanner"
)

type status int

const (
	initial status = iota
	annotationName
	attributeName
	attributeValue
	done
)

func parseAnnotation(line string) (Annotation, error) {
	withoutComment := strings.TrimLeft(strings.TrimSpace(line), "/")

	annotation := Annotation{
		Name:       "",
		Attributes: make(map[string]string),
	}

	var s scanner.Scanner
	s.Init(strings.NewReader(withoutComment))

	var tok rune
	currentStatus := initial
	var attrName string

	for tok != scanner.EOF && currentStatus < done {
		tok = s.Scan()
		switch tok {
		case '@':
			currentStatus = annotationName
		case '(':
			currentStatus = attributeName
		case '=':
			currentStatus = attributeValue
		case ',':
			currentStatus = attributeName
		case ')':
			currentStatus = done
		case scanner.Ident:
			switch currentStatus {
			case annotationName:
				annotation.Name = s.TokenText()
			case attributeName:
				attrName = s.TokenText()
			}
		default:
			switch currentStatus {
			case attributeValue:
				annotation.Attributes[strings.ToLower(attrName)] = strings.Trim(s.TokenText(), "\"")
			}
		}
	}

	if currentStatus != done {
		return annotation, fmt.Errorf("Invalid completion-status %v for annotation:%s", currentStatus, line)
	}
	return annotation, nil
}
