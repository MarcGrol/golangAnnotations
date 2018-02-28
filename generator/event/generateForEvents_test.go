package event

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/MarcGrol/golangAnnotations/generator/filegen"
	"github.com/MarcGrol/golangAnnotations/model"
	"github.com/stretchr/testify/assert"
)

func cleanup() {
	os.Remove(filegen.Prefixed("./testData/ast.json"))
	os.Remove(filegen.Prefixed("./testData/aggregates.go"))
	os.Remove(filegen.Prefixed("./testData/wrappers.go"))
	os.Remove(filegen.Prefixed("./store/testDataStore/testDataEventStore.go"))
}

func TestGenerateForEvents(t *testing.T) {
	cleanup()
	defer cleanup()

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
	err := NewGenerator().Generate("testData", model.ParsedSources{Structs: s})
	assert.Nil(t, err)

	// check that generated files exisst
	_, err = os.Stat(filegen.Prefixed("./testData/aggregates.go"))
	assert.NoError(t, err)

	data, err := ioutil.ReadFile(filegen.Prefixed("./testData/aggregates.go"))
	assert.NoError(t, err)
	assert.Contains(t, string(data), "type TestAggregate interface {")
	assert.Contains(t, string(data), "ApplyMyStruct(c context.Context, event MyStruct)")
	assert.Contains(t, string(data), "func ApplyTestEvent(c context.Context, envelope envelope.Envelope, aggregateRoot TestAggregate) error {")
	assert.Contains(t, string(data), "func ApplyTestEvents(c context.Context, envelopes []envelope.Envelope, aggregateRoot TestAggregate) error {")
	assert.Contains(t, string(data), "func UnWrapTestEvent(envelope *envelope.Envelope) (envelope.Event, error) {")
	assert.Contains(t, string(data), "func UnWrapTestEvents(envelopes []envelope.Envelope) ([]envelope.Event, error) {")

	// check that generate code has 4 helper functions for MyStruct
	data, err = ioutil.ReadFile(filegen.Prefixed("./testData/wrappers.go"))
	assert.NoError(t, err)
	assert.Contains(t, string(data), "func (s *MyStruct) Wrap(rc request.Context) (*envelope.Envelope,error) {")
	assert.Contains(t, string(data), "func IsMyStruct(envelope *envelope.Envelope) bool {")
	assert.Contains(t, string(data), "func GetIfIsMyStruct(envelope *envelope.Envelope) (*MyStruct, bool) {")
	assert.Contains(t, string(data), "func UnWrapMyStruct(envelope *envelope.Envelope) (*MyStruct,error) {")

	_, err = os.Stat(filegen.Prefixed("./testData/wrappers.go"))
	assert.NoError(t, err)

	os.Remove(filegen.Prefixed("./testData/ast.json"))
	os.Remove(filegen.Prefixed("./testData/aggregates.go"))
	os.Remove(filegen.Prefixed("./testData/wrappers.go"))
	os.Remove(filegen.Prefixed("./testData/wrappers_test.go"))
	os.Remove(filegen.Prefixed("./repository/storeEvents.go"))

}

func TestIsEvent(t *testing.T) {
	s := model.Struct{
		DocLines: []string{
			`//@Event( aggregate = "person")`},
	}
	assert.True(t, isEvent(s))
}

func TestGetAggregateName(t *testing.T) {
	s := model.Struct{
		DocLines: []string{
			`//@Event( aggregate = "person")`},
	}
	assert.Equal(t, "person", getAggregateName(s))
}
