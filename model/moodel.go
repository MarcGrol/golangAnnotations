package model

import (
	"fmt"
	"strings"
)

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

func (s Struct) getEventAggregateAnnotation() (string, bool) {
	found := false
	aggregateName := ""

	for _, line := range s.DocLines {
		count, err := fmt.Sscanf(strings.TrimSpace(line), "// +event -> aggregate: %s",
			&aggregateName)
		if err == nil && count == 1 {
			found = true
			break
		}
	}
	return aggregateName, found

}
