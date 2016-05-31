package generator

import (
	"testing"

	"github.com/MarcGrol/astTools/model"
	"github.com/stretchr/testify/assert"
)

func TestGenerateForStructs(t *testing.T) {
	s := []model.Struct{
		{
			DocLines:    []string{"// +event -> aggregate: Person"},
			PackageName: "generator",
			Name:        "MyStruct",
			Fields: []model.Field{
				{Name: "StringField", TypeName: "string", IsPointer: false, IsSlice: false},
				{Name: "IntField", TypeName: "int", IsPointer: false, IsSlice: false},
				{Name: "StructField", TypeName: "MyStruct", IsPointer: true, IsSlice: false},
				{Name: "SliceField", TypeName: "MyStruct", IsPointer: false, IsSlice: true},
			},
		},
	}
	err := GenerateForStructs(".", s)
	assert.Nil(t, err)

}
