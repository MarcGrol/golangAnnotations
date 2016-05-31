package model

import (
	"fmt"
	"strings"
)

type Service struct {
	DocLines     []string
	PackageName  string
	Name         string
	Methods      []Method
	CommentLines []string
}

type Method struct {
	DocLines     []string
	PackageName  string
	Service      *Service
	Name         string
	InputArgs    []Struct
	OutputArgs   []Struct
	CommentLines []string
}

type Struct struct {
	DocLines     []string
	PackageName  string
	Name         string
	Fields       []Field
	CommentLines []string
}

type Field struct {
	DocLines     []string
	Name         string
	TypeName     string
	IsSlice      bool
	IsPointer    bool
	Tag          string
	CommentLines []string
}

func (s Struct) IsEvent() bool {
	_, hasEventAnnotation := s.getEventAggregateAnnotation()
	return hasEventAnnotation
}

func (s Struct) GetAggregateName() string {
	aggr, _ := s.getEventAggregateAnnotation()
	return aggr
}

func (s Struct) IsRestService() bool {
	return false
}

func (s Struct) GetRestServiceParamaters() (path string) {
	return ""
}

func (m Method) IsRestMethod() bool {
	return false
}

func (s Struct) GetRestMethodParamaters() (path string, method string) {
	return "", "GET"
}

func (s Struct) getEventAggregateAnnotation() (string, bool) {
	found := false
	aggregateName := ""

	for _, line := range s.DocLines {
		//log.Printf("line:%s", line)
		count, err := fmt.Sscanf(strings.TrimSpace(line), "// +event -> aggregate: %s",
			&aggregateName)
		if err == nil && count == 1 {
			//	log.Printf("Match:%s", line)
			found = true
			break
		} else {
			//	log.Printf("No match:%s", line)
		}
	}
	return aggregateName, found

}
