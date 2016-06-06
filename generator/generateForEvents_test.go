package generator

import (
	"os"
	"strings"
	"testing"

	"io/ioutil"

	"github.com/MarcGrol/astTools/model"
	"github.com/stretchr/testify/assert"
)

func TestGenerateForStructs(t *testing.T) {
	os.Remove("./testData/aggregates.go")
	os.Remove("./testData/wrappers.go")

	s := []model.Struct{
		{
			DocLines:    []string{`// {"Annotation":"Event","With":{"Aggregate":"Test"}}`},
			PackageName: "testData",
			Name:        "MyStruct",
			Fields: []model.Field{
				{Name: "StringField", TypeName: "string", IsPointer: false, IsSlice: false},
				{Name: "IntField", TypeName: "int", IsPointer: false, IsSlice: false},
				{Name: "StructField", TypeName: "MyStruct", IsPointer: true, IsSlice: false},
				{Name: "SliceField", TypeName: "MyStruct", IsPointer: false, IsSlice: true},
			},
		},
	}
	err := GenerateForStructs("testData", s)
	assert.Nil(t, err)

	// check that generated files exisst
	_, err = os.Stat("./testData/aggregates.go")
	assert.NoError(t, err)

	// check that generate code has 4 helper functions for MyStruct
	data, err := ioutil.ReadFile("./testData/wrappers.go")
	assert.NoError(t, err)
	assert.True(t, strings.Contains(string(data), "func (s *MyStruct) Wrap(uid string) (*Envelope,error) {"))
	assert.True(t, strings.Contains(string(data), "func IsMyStruct(envelope *Envelope) bool {"))
	assert.True(t, strings.Contains(string(data), "func GetIfIsMyStruct(envelop *Envelope) (*MyStruct, bool) {"))
	assert.True(t, strings.Contains(string(data), "func UnWrapMyStruct(envelop *Envelope) (*MyStruct,error) {"))

	_, err = os.Stat("./testData/wrappers.go")
	assert.NoError(t, err)
	data, err = ioutil.ReadFile("./testData/aggregates.go")
	assert.NoError(t, err)
	t.Log("data:" + string(data))
	assert.True(t, strings.Contains(string(data), "type TestAggregate interface {"))
	assert.True(t, strings.Contains(string(data), "ApplyMyStruct(event MyStruct)"))
	assert.True(t, strings.Contains(string(data), "func ApplyTestEvent(envelop Envelope, aggregateRoot TestAggregate) error {"))
	assert.True(t, strings.Contains(string(data), "func ApplyTestEvents(envelopes []Envelope, aggregateRoot TestAggregate) error {"))
}
