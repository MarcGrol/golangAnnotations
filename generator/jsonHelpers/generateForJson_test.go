package jsonHelpers

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/MarcGrol/golangAnnotations/model"
	"github.com/stretchr/testify/assert"
)

func cleanup() {
	os.Remove("./testData/»ast.json")
	os.Remove("./testData/»example_json.go")
}

func TestGenerateForJson(t *testing.T) {
	cleanup()
	defer cleanup()

	e := []model.Enum{
		{
			PackageName: "testData",
			Filename:    "example.go",
			DocLines:    []string{"// @JsonEnum()"},
			Name:        "ColorType",
			EnumLiterals: []model.EnumLiteral{
				{Name: "ColorTypeRed"},
				{Name: "ColorTypeGreen"},
				{Name: "ColorTypeBlue"},
			},
		},
	}

	s := []model.Struct{
		{
			PackageName: "testData",
			Filename:    "example.go",
			DocLines:    []string{`// @JsonStruct()`},
			Name:        "ColoredThing",
			Fields: []model.Field{
				{
					Name:     "Name",
					TypeName: "string",
				},
				{
					Name:     "Tags",
					TypeName: "string",
					IsSlice:  true,
				},
				{
					Name:     "PrimaryColor",
					TypeName: "ColorType",
				},
				{
					Name:     "OtherColors",
					TypeName: "ColorType",
					IsSlice:  true,
				},
			},
		},
	}

	ps := model.ParsedSources{
		Enums:   e,
		Structs: s,
	}
	err := NewGenerator().Generate("./testData/", ps)
	assert.Nil(t, err)

	// check that generated files exisst
	_, err = os.Stat("./testData/»example_json.go")
	assert.NoError(t, err)

	// check that generate code has 4 helper functions for MyStruct
	data, err := ioutil.ReadFile("./testData/»example_json.go")
	assert.NoError(t, err)
	assert.Contains(t, string(data), `func (r *ColorType) UnmarshalJSON(data []byte) error {`)
	assert.Contains(t, string(data), `func (r ColorType) MarshalJSON() ([]byte, error) {`)

	assert.Contains(t, string(data), `func (data *ColoredThing) UnmarshalJSON(b []byte) error {`)
	assert.Contains(t, string(data), `func (data ColoredThing) MarshalJSON() ([]byte, error) {`)

}

func TestIsJsonEnum(t *testing.T) {
	e := model.Enum{
		DocLines: []string{
			`// @JsonStruct()`,
			`// @JsonEnum()`,
		},
	}
	assert.True(t, IsJSONEnum(e))
}

func TestIsJsonStruct(t *testing.T) {
	s := model.Struct{
		DocLines: []string{
			`// @Event(aggregate = "Test")`,
			`// @JsonStruct()`,
		},
	}
	assert.True(t, IsJSONStruct(s))
}
