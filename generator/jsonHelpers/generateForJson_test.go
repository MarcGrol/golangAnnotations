package jsonHelpers

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/MarcGrol/golangAnnotations/generator/event"
	"github.com/MarcGrol/golangAnnotations/generator/event/eventAnnotation"
	"github.com/MarcGrol/golangAnnotations/generator/jsonHelpers/jsonAnnotation"
	"github.com/MarcGrol/golangAnnotations/model"
	"github.com/stretchr/testify/assert"
)

func cleanup() {
	os.Remove("./testData/jsonHelpers.go")
}

func TestGenerateForJson(t *testing.T) {
	cleanup()
	defer cleanup()

	e := []model.Enum{
		{
			DocLines:    []string{"// @JsonEnum()"},
			PackageName: "testData",
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
			DocLines:    []string{`// @JsonStruct()`},
			PackageName: "testData",
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
	err := Generate("testData", ps)
	assert.Nil(t, err)

	// check that generated files exisst
	_, err = os.Stat("./testData/jsonHelpers.go")
	assert.NoError(t, err)

	// check that generate code has 4 helper functions for MyStruct
	data, err := ioutil.ReadFile("./testData/jsonHelpers.go")
	assert.NoError(t, err)
	assert.Contains(t, string(data), `func (r *ColorType) UnmarshalJSON(data []byte) error {`)
	assert.Contains(t, string(data), `func (r ColorType) MarshalJSON() ([]byte, error) {`)

	assert.Contains(t, string(data), `func (data *ColoredThing) UnmarshalJSON(b []byte) error {`)
	assert.Contains(t, string(data), `func (data ColoredThing) MarshalJSON() ([]byte, error) {`)

}

func TestIsJsonEnum(t *testing.T) {
	jsonAnnotation.Register()
	e := model.Enum{
		DocLines: []string{
			`// @JsonStruct()`,
			`// @JsonEnum()`,
		},
	}
	assert.True(t, IsJsonEnum(e))
}

func TestIsJsonStruct(t *testing.T) {
	eventAnnotation.Register()
	jsonAnnotation.Register()
	s := model.Struct{
		DocLines: []string{
			`// @Event(aggregate = "Test")`,
			`// @JsonStruct()`,
		},
	}
	assert.True(t, event.IsEvent(s))
	assert.True(t, IsJsonStruct(s))
}
