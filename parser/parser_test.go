package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseStructsInFile(t *testing.T) {

	structs, err := FindStructsInFile("testData/example.go")
	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(structs))

	{
		s := structs[0]
		assert.Equal(t, "Person", s.Name)
		assert.Equal(t, 9, len(s.Fields))

		assertField(t,
			Field{Name: "FirstName", TypeName: "string", IsPointer: false, IsSlice: false},
			s.Fields[0])

		assertField(t,
			Field{Name: "LastName", TypeName: "string", IsPointer: false, IsSlice: false},
			s.Fields[1])

		assertField(t,
			Field{Name: "Age", TypeName: "int", IsPointer: false, IsSlice: false},
			s.Fields[2])

		assertField(t,
			Field{Name: "Nice", TypeName: "bool", IsPointer: true, IsSlice: false},
			s.Fields[3])

		assertField(t,
			Field{Name: "Color", TypeName: "ColorType", IsPointer: false, IsSlice: false},
			s.Fields[4])

		assertField(t,
			Field{Name: "OptionalColor", TypeName: "ColorType", IsPointer: true, IsSlice: false},
			s.Fields[5])

		assertField(t,
			Field{Name: "Father", TypeName: "Person", IsPointer: true, IsSlice: false},
			s.Fields[6])

		assertField(t,
			Field{Name: "Uncles", TypeName: "Person", IsPointer: true, IsSlice: true},
			s.Fields[7])

		assertField(t,
			Field{Name: "Children", TypeName: "Person", IsPointer: false, IsSlice: true},
			s.Fields[8])

	}
}

func TestParseStructsInDir(t *testing.T) {

	structs, err := FindStructsInDir("testData", ".*xample.*")
	assert.Equal(t, nil, err)
	assert.Equal(t, 3, len(structs))

	// Order is undetermined
	for _, s := range structs {
		if s.Name == "Person" {
			assert.Equal(t, 9, len(structs[0].Fields))
		}
		if s.Name == "MyStruct" {
			assert.Equal(t, 1, len(structs[1].Fields))
		}
		if s.Name == "OtherStruct" {
			assert.Equal(t, 1, len(structs[2].Fields))
		}
	}
}

func assertField(t *testing.T, expected Field, actual Field) {
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.TypeName, actual.TypeName)
	assert.Equal(t, expected.IsPointer, actual.IsPointer)
	assert.Equal(t, expected.IsSlice, actual.IsSlice)
}
