package testData

import "fmt"

type ColorType int

const (
	Green ColorType = iota
	Yellow
	Red
)

func MyFunc(value int, yes bool) (string, error) {
	return fmt.Sprintf("%d", value), nil
}

// Struct comment before type
type Person struct {
	FirstName, LastName string
	Age                 int // Age comment
	// Before nice comment
	Nice *bool // After Nice comment
	// Before Color comment
	Color         ColorType `json:"COLOR_TYPE"`
	OptionalColor *ColorType
	Father        *Person
	Uncles        []*Person
	Children      []Person
}

func (p Person) Dump() string {
	return fmt.Sprintf("%+v", p)
}

type MyStruct struct {
	X int
}
