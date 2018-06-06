package structs

import "fmt"

// Struct comment before type
type Person struct {
	FirstName, LastName string
	Age                 int // Age comment
	// Before nice comment
	Nice *bool // After Nice comment
	// Before Color comment
	Color              ColorType `json:"COLOR_TYPE"`
	OptionalColor      *ColorType
	Father             *Person
	Uncles             []*Person
	Children           []Person
	ChildMap           map[string]Person
	ChildPointerMap    map[string]*Person
	ChildrenMap        map[string][]Person
	ChildrenPointerMap map[string][]*Person
}

type ColorType int

const (
	Green ColorType = iota
	Yellow
	Red
)

func MyFunc(value int, yes bool) (string, error) {
	return fmt.Sprintf("%d", value), nil
}

func (p Person) Dump() string {
	return fmt.Sprintf("%+v", p)
}

type MyStruct struct {
	X int
}
