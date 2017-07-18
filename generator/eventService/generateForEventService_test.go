package eventService

import (
	"fmt"
	"os"
	"testing"

	"io/ioutil"

	"github.com/MarcGrol/golangAnnotations/generator/rest/restAnnotation"
	"github.com/MarcGrol/golangAnnotations/model"
	"github.com/stretchr/testify/assert"
)

func cleanup() {
	os.Remove("./testData/$httpMyEventService.go")
	os.Remove("./testData/$eventHandler.go")
}

func TestGenerateForWeb(t *testing.T) {
	cleanup()
	defer cleanup()

	s := []model.Struct{
		{
			DocLines:    []string{`// @EventService( self = "self", subscriptions = "other")`},
			PackageName: "testData",
			Name:        "MyEventService",
			Operations:  []*model.Operation{},
		},
	}

	s[0].Operations = append(s[0].Operations,
		&model.Operation{
			DocLines:      []string{`// @EventOperation()`},
			Name:          "doit",
			RelatedStruct: &model.Field{TypeName: "MyService"},
			InputArgs: []model.Field{
				{Name: "c", TypeName: "context.Context"},
				{Name: "structExample", TypeName: "events.OrderCreated"},
			},
			OutputArgs: []model.Field{
				{TypeName: "error"},
			},
		})

	err := Generate("testData", model.ParsedSources{Structs: s})
	assert.Nil(t, err)

	// check that generated files exisst
	_, err = os.Stat("./testData/$eventHandler.go")
	assert.NoError(t, err)

	// check that generate code has 4 helper functions for MyStruct
	data, err := ioutil.ReadFile("./testData/$eventHandler.go")
	assert.NoError(t, err)
	assert.Contains(t, string(data), `bus.Subscribe("other", subscriber, es.handleEvent)`)
	assert.Contains(t, string(data), `func (es *MyEventService) handleEvent(c context.Context, topic string, envelope events.Envelope) {`)
}

func TestIsRestService(t *testing.T) {
	restAnnotation.Register()
	s := model.Struct{
		DocLines: []string{
			`//@EventService( self = "me", subscriptions = "other")`},
	}
	assert.True(t, IsEventService(s))
}

func TestGetEventServiceSelf(t *testing.T) {
	restAnnotation.Register()
	s := model.Struct{
		DocLines: []string{
			`//@EventService( self = "me", subscriptions = "other1, other2")`},
	}
	assert.Equal(t, "me", GetEventServiceSelfName(s))
	assert.Equal(t, []string{"other1", "other2"}, GetEventServiceSubscriptions(s))
}

func TestIsEventOperation(t *testing.T) {
	assert.True(t, IsEventOperation(createOper()))
}

func TestGetEventName(t *testing.T) {
	assert.Equal(t, "OrderCreated", GetInputArgType(createOper()))
}

func TestGetInputArgTypePerson(t *testing.T) {
	assert.Equal(t, "OrderCreated", GetInputArgType(createOper()))
}

func createOper() model.Operation {
	restAnnotation.Register()
	o := model.Operation{
		DocLines: []string{
			fmt.Sprintf("//@EventOperation()"),
		},
		InputArgs: []model.Field{
			{Name: "ctx", TypeName: "context.Context"},
			{Name: "uid", TypeName: "events.OrderCreated"},
		},
		OutputArgs: []model.Field{
			{TypeName: "error"},
		},
	}
	return o
}
