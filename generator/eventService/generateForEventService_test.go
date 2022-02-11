package eventService

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/f0rt/golangAnnotations/generator/generationUtil"
	"github.com/f0rt/golangAnnotations/model"
	"github.com/stretchr/testify/assert"
)

func cleanup() {
	os.Remove(generationUtil.Prefixed("./testData/ast.json"))
	os.Remove(generationUtil.Prefixed("./testData/httpMyEventService.go"))
	os.Remove(generationUtil.Prefixed("./testData/eventHandler.go"))
	os.Remove(generationUtil.Prefixed("./testData/eventHandlerHelpers_test.go"))
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

	err := NewGenerator().Generate("testData", model.ParsedSources{Structs: s})
	assert.Nil(t, err)

	// check that generated files exisst
	_, err = os.Stat(generationUtil.Prefixed("./testData/eventHandler.go"))
	assert.NoError(t, err)

	// check that generate code has 4 helper functions for MyStruct
	data, err := ioutil.ReadFile(generationUtil.Prefixed("./testData/eventHandler.go"))
	assert.NoError(t, err)
	assert.Contains(t, string(data), `bus.Subscribe("other", subscriber, es.handleOrEnqueueEvent)`)
	assert.Contains(t, string(data), `func (es *MyEventService) handleEvent(c context.Context, rc request.Context, topic string, envlp envelope.Envelope) error {`)
}

func TestIsRestService(t *testing.T) {
	s := model.Struct{
		DocLines: []string{
			`//@EventService( self = "me")`},
	}
	assert.True(t, IsEventService(s))
}

func TestGetEventServiceSelf(t *testing.T) {
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
