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

// Struct Person todo
type Person struct {
	FirstName, LastName string // FirstName
	Age                 int    // Age
	Nice                *bool  // Nice
	Color               ColorType
	OptionalColor       *ColorType
	Father              *Person
	Uncles              []*Person
	Children            []Person
}

func (p Person) Dump() string {
	return fmt.Sprintf("%+v", p)
}

type MyStruct struct {
	X int
}
