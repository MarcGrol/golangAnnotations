package rest

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
	os.Remove("./testData/$httpMyService.go")
	os.Remove("./testData/$httpMyServiceHelpers_test.go")
	os.Remove("./testData/$httpClientForMyService.go")
	os.Remove("./testData/$httpMyServiceHelpers_test.go")
}

func TestGenerateForWeb(t *testing.T) {
	cleanup()
	defer cleanup()

	s := []model.Struct{
		{
			DocLines:    []string{"// @RestService( path = \"/api\")"},
			PackageName: "testData",
			Name:        "MyService",
			Operations:  []*model.Operation{},
		},
	}

	s[0].Operations = append(s[0].Operations,
		&model.Operation{
			DocLines:      []string{"// @RestOperation(path = \"/person\", method = \"GET\", format = \"JSON\", form = \"true\" )"},
			Name:          "doit",
			RelatedStruct: &model.Field{TypeName: "MyService"},
			InputArgs: []model.Field{
				{Name: "uid", TypeName: "int"},
				{Name: "subuid", TypeName: "string"},
			},
			OutputArgs: []model.Field{
				{TypeName: "error"},
			},
		})

	err := Generate("testData", model.ParsedSources{Structs: s})
	assert.Nil(t, err)

	{
		// check that generated files exists
		_, err = os.Stat("./testData/$httpMyService.go")
		assert.NoError(t, err)

		// check that generate code has 4 helper functions for MyStruct
		data, err := ioutil.ReadFile("./testData/$httpMyService.go")
		assert.NoError(t, err)
		assert.Contains(t, string(data), "func (ts *MyService) HTTPHandler() http.Handler {")
		assert.Contains(t, string(data), "func doit( service *MyService ) http.HandlerFunc {")

	}
	{
		// check that generated files exists
		_, err = os.Stat("./testData/$httpMyService.go")
		assert.NoError(t, err)

		// check that generate code has 4 helper functions for MyStruct
		data, err := ioutil.ReadFile("./testData/$httpMyServiceHelpers_test.go")
		assert.NoError(t, err)
		assert.Contains(t, string(data), "func doitTestHelper")
	}

	{
		// check that generated files exists
		_, err = os.Stat("./testData/$httpClientForMyService.go")
		assert.NoError(t, err)

		// check that generate code has 4 helper functions for MyStruct
		data, err := ioutil.ReadFile("./testData/$httpClientForMyService.go")
		assert.NoError(t, err)
		assert.Contains(t, string(data), "func (c *HTTPClient) Doit(ctx context.Context, url string , cookie *http.Cookie, requestUID string, timeout time.Duration)  (int ,*errorh.Error,error) {")
	}

}

func TestIsRestService(t *testing.T) {
	restAnnotation.Register()
	s := model.Struct{
		DocLines: []string{
			`//@RestService( path = "/api")`},
	}
	assert.True(t, isRestService(s))
}

func TestGetRestServicePath(t *testing.T) {
	restAnnotation.Register()
	s := model.Struct{
		DocLines: []string{
			`//@RestService( path = "/api")`},
	}
	assert.Equal(t, "/api", getRestServicePath(s))
}

func TestIsRestOperation(t *testing.T) {
	assert.True(t, isRestOperation(createOper("GET")))
}

func TestGetRestOperationMethod(t *testing.T) {
	assert.Equal(t, "GET", getRestOperationMethod(createOper("GET")))
}

func TestGetRestOperationPath(t *testing.T) {
	assert.Equal(t, "/api/person", getRestOperationPath(createOper("DONTCARE")))
}

func TestHasInputGet(t *testing.T) {
	assert.False(t, hasInput(createOper("GET")))
}

func TestHasInputDelete(t *testing.T) {
	assert.False(t, hasInput(createOper("DELETE")))
}

func TestHasInputPost(t *testing.T) {
	assert.True(t, hasInput(createOper("POST")))
}

func TestHasInputPut(t *testing.T) {
	assert.True(t, hasInput(createOper("PUT")))
}

func TestGetInputArgTypeString(t *testing.T) {
	restAnnotation.Register()
	o := model.Operation{
		InputArgs: []model.Field{
			{TypeName: "string"},
		},
	}
	assert.Equal(t, "", getInputArgType(o))
}

func TestGetInputArgTypePerson(t *testing.T) {
	assert.Equal(t, "Person", getInputArgType(createOper("DONTCARE")))
}

func TestGetInputArgName(t *testing.T) {
	assert.Equal(t, "person", getInputArgName(createOper("DONTCARE")))
}

func TestGetInputParamString(t *testing.T) {
	assert.Equal(t, "ctx,uid,person", getInputParamString(createOper("DONTCARE")))
}

func TestHasOutput(t *testing.T) {
	assert.True(t, hasOutput(createOper("DONTCARE")))
}

func TestGetOutputArgType(t *testing.T) {
	assert.Equal(t, "Person", getOutputArgType(createOper("DONTCARE")))
}

func TestIsPrimitiveTrue(t *testing.T) {
	f := model.Field{Name: "uid", TypeName: "string"}
	assert.True(t, isPrimitiveArg(f))
}

func TestIsPrimitiveFalse(t *testing.T) {
	f := model.Field{Name: "person", TypeName: "Person"}
	assert.False(t, isPrimitiveArg(f))
}

func TestIsNumberTrue(t *testing.T) {
	f := model.Field{Name: "uid", TypeName: "int"}
	assert.True(t, isNumberArg(f))
}

func TestIsNumberFalse(t *testing.T) {
	f := model.Field{Name: "uid", TypeName: "string"}
	assert.False(t, isNumberArg(f))
}

func createOper(method string) model.Operation {
	restAnnotation.Register()
	o := model.Operation{
		DocLines: []string{
			fmt.Sprintf("//@RestOperation( method = \"%s\", path = \"/api/person\")", method),
		},
		InputArgs: []model.Field{
			{Name: "ctx", TypeName: "context.Context"},
			{Name: "uid", TypeName: "string"},
			{Name: "person", TypeName: "Person"},
		},
		OutputArgs: []model.Field{
			{TypeName: "Person"},
			{TypeName: "error"},
		},
	}
	return o
}
