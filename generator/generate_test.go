package generator

import (
	"testing"

	"github.com/MarcGrol/astTools/parser"
	"github.com/stretchr/testify/assert"
)

func TestGenerateForStruct(t *testing.T) {

	err := GenerateForStruct(
		parser.Struct{
			PackageName: "generator",
			Name:        "MyStruct",
			Fields: []parser.Field{
				{Name: "StringField", TypeName: "string", IsPointer: false, IsSlice: false},
				{Name: "IntField", TypeName: "int", IsPointer: false, IsSlice: false},
				{Name: "StructField", TypeName: "MyStruct", IsPointer: true, IsSlice: false},
				{Name: "SliceField", TypeName: "MyStruct", IsPointer: false, IsSlice: true},
			},
		})
	assert.Nil(t, err)

}

func TestParseGenerate(t *testing.T) {
	//structs, err := parser.FindStructsInFile("testInput.go")
}
