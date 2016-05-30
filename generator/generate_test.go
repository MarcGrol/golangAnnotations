package generator

import (
	"testing"

	"github.com/MarcGrol/astTools/model"
	"github.com/stretchr/testify/assert"
)

func TestGenerateForStruct(t *testing.T) {

	s := model.Struct{
		DocLines:    []string{"// +event -> aggregate: person"},
		PackageName: "generator",
		Name:        "MyStruct",
		Fields: []model.Field{
			{Name: "StringField", TypeName: "string", IsPointer: false, IsSlice: false},
			{Name: "IntField", TypeName: "int", IsPointer: false, IsSlice: false},
			{Name: "StructField", TypeName: "MyStruct", IsPointer: true, IsSlice: false},
			{Name: "SliceField", TypeName: "MyStruct", IsPointer: false, IsSlice: true},
		},
	}
	t.Logf("struct to generate:%+v", s)
	err := GenerateForStruct(s)
	assert.Nil(t, err)

}

func TestParseGenerate(t *testing.T) {
	//structs, err := parser.FindStructsInFile("testInput.go")
}
