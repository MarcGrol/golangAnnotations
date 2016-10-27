package event

import (
	"os"
	"testing"

	"io/ioutil"

	"github.com/MarcGrol/golangAnnotations/generator/event/eventAnnotation"
	"github.com/MarcGrol/golangAnnotations/model"
	"github.com/stretchr/testify/assert"
)

func TestGenerateForEvents(t *testing.T) {
	os.Remove("./testData/aggregates.go")
	os.Remove("./testData/wrappers.go")

	s := []model.Struct{
		{
			PackageName: "testData",
			DocLines:    []string{`//@Event(aggregate = "Test")`},
			Name:        "MyStruct",
			Fields: []model.Field{
				{Name: "StringField", TypeName: "string", IsPointer: false, IsSlice: false},
				{Name: "IntField", TypeName: "int", IsPointer: false, IsSlice: false},
				{Name: "StructField", TypeName: "MyStruct", IsPointer: true, IsSlice: false},
				{Name: "SliceField", TypeName: "MyStruct", IsPointer: false, IsSlice: true},
			},
		},
	}
	err := Generate("testData", model.ParsedSources{Structs: s})
	assert.Nil(t, err)

	// check that generated files exisst
	_, err = os.Stat("./testData/aggregates.go")
	assert.NoError(t, err)

	// check that generate code has 4 helper functions for MyStruct
	data, err := ioutil.ReadFile("./testData/wrappers.go")
	assert.NoError(t, err)
	assert.Contains(t, string(data), "func (s *MyStruct) Wrap(sessionUID string) (*Envelope,error) {")
	assert.Contains(t, string(data), "func IsMyStruct(envelope *Envelope) bool {")
	assert.Contains(t, string(data), "func GetIfIsMyStruct(envelop *Envelope) (*MyStruct, bool) {")
	assert.Contains(t, string(data), "func UnWrapMyStruct(envelop *Envelope) (*MyStruct,error) {")

	_, err = os.Stat("./testData/wrappers.go")
	assert.NoError(t, err)
	data, err = ioutil.ReadFile("./testData/aggregates.go")
	assert.NoError(t, err)
	assert.Contains(t, string(data), "type TestAggregate interface {")
	assert.Contains(t, string(data), "ApplyMyStruct(c context.Context, event MyStruct)")
	assert.Contains(t, string(data), "func ApplyTestEvent(c context.Context, envelop Envelope, aggregateRoot TestAggregate) error {")
	assert.Contains(t, string(data), "func ApplyTestEvents(c context.Context, envelopes []Envelope, aggregateRoot TestAggregate) error {")

	os.Remove("./testData/aggregates.go")
	os.Remove("./testData/wrappers.go")
	os.Remove("./testData/wrappers_test.go")
	os.Remove("./repo/storeEvents.go")

}

func TestIsEvent(t *testing.T) {
	eventAnnotation.Register()
	s := model.Struct{
		DocLines: []string{
			`//@Event( aggregate = "person")`},
	}
	assert.True(t, IsEvent(s))
}

func TestGetAggregateName(t *testing.T) {
	eventAnnotation.Register()
	s := model.Struct{
		DocLines: []string{
			`//@Event( aggregate = "person")`},
	}
	assert.Equal(t, "person", GetAggregateName(s))
}
