package eventService

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

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
			DocLines:    []string{`// @EventService( self = "self", async="true", admin="true" )`},
			PackageName: "testData",
			Name:        "MyEventService",
			Operations: []*model.Operation{
				{
					DocLines:      []string{`// @EventOperation( topic = "other" )`},
					Name:          "doit",
					RelatedStruct: &model.Field{TypeName: "MyService"},
					InputArgs: []model.Field{
						{Name: "c", TypeName: "context.Context"},
						{Name: "structExample", TypeName: "events.OrderCreated"},
					},
					OutputArgs: []model.Field{
						{TypeName: "error"},
					},
				},
			},
		},
	}

	err := Generate("testData", model.ParsedSources{Structs: s})
	assert.Nil(t, err)

	// check that generated files exisst
	_, err = os.Stat("./testData/$eventHandler.go")
	assert.NoError(t, err)

	// check that generate code has 4 helper functions for MyStruct
	data, err := ioutil.ReadFile("./testData/$eventHandler.go")
	assert.NoError(t, err)
	assert.Contains(t, string(data), `bus.Subscribe("other", subscriber, es.handleEvent)`)
	assert.Contains(t, string(data), `func (es *MyEventService) handleEvent(c context.Context, credentials rest.Credentials, topic string, envelope envelope.Envelope) {`)
}

func TestIsRestService(t *testing.T) {
	restAnnotation.Register()
	s := model.Struct{
		DocLines: []string{
			`//@EventService( self = "me")`},
	}
	assert.True(t, IsEventService(s))
}

func TestGetEventServiceSelf(t *testing.T) {
	restAnnotation.Register()
	s := model.Struct{
		DocLines: []string{
			`//@EventService( self = "me" )`},
	}
	assert.Equal(t, "me", GetEventServiceSelfName(s))
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
			fmt.Sprintf("//@EventOperation( topic = \"other1\" )"),
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
