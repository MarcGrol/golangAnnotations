package rest

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/MarcGrol/golangAnnotations/generator/generationUtil"
	"github.com/MarcGrol/golangAnnotations/model"
	"github.com/stretchr/testify/assert"
)

func cleanup() {
	os.Remove(generationUtil.Prefixed("./testData/ast.json"))
	os.Remove(generationUtil.Prefixed("./testData/httpMyService.go"))
	os.Remove(generationUtil.Prefixed("./testData/httpMyServiceHelpers_test.go"))
	os.Remove(generationUtil.Prefixed("./testData/httpClientForMyService.go"))
	os.Remove(generationUtil.Prefixed("./testData/httpMyServiceHelpers_test.go"))
	os.Remove(generationUtil.Prefixed("./testData/testDataTestLog/httpTestMyService.go"))
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
	{
		err := NewGenerator().Generate("testData", model.ParsedSources{Structs: s})
		assert.Nil(t, err)
	}

	{
		{
			// check that generated files exists
			_, err := os.Stat(generationUtil.Prefixed("./testData/httpMyService.go"))
			assert.NoError(t, err)
		}
		{
			// check that generate code has 4 helper functions for MyStruct
			data, err := ioutil.ReadFile(generationUtil.Prefixed("./testData/httpMyService.go"))
			assert.NoError(t, err)
			assert.Contains(t, string(data), "func (ts *MyService) HTTPHandler() http.Handler {")
			assert.Contains(t, string(data), "func doit( service *MyService ) http.HandlerFunc {")
		}
	}
	{
		{
			// check that generated files exists
			_, err := os.Stat(generationUtil.Prefixed("./testData/httpMyService.go"))
			assert.NoError(t, err)
		}
		{
			// check that generate code has 4 helper functions for MyStruct
			data, err := ioutil.ReadFile(generationUtil.Prefixed("./testData/httpMyServiceHelpers_test.go"))
			assert.NoError(t, err)
			assert.Contains(t, string(data), "func doitTestHelper")
		}
	}

	{
		{
			// check that generated files exists
			_, err := os.Stat(generationUtil.Prefixed("./testData/httpClientForMyService.go"))
			assert.NoError(t, err)
		}
		{
			// check that generate code has 4 helper functions for MyStruct
			data, err := ioutil.ReadFile(generationUtil.Prefixed("./testData/httpClientForMyService.go"))
			assert.NoError(t, err)
			assert.Contains(t, string(data), "func (c *HTTPClient) Doit(ctx context.Context, url string , cookie *http.Cookie, requestUID string, timeout time.Duration)  (int ,*errorh.Error,error) {")
		}
	}

}

func TestIsRestService(t *testing.T) {
	s := model.Struct{
		DocLines: []string{
			`//@RestService( path = "/api")`},
	}
	assert.True(t, IsRestService(s))
}

func TestGetRestServicePath(t *testing.T) {
	s := model.Struct{
		DocLines: []string{
			`//@RestService( path = "/api")`},
	}
	assert.Equal(t, "/api", GetRestServicePath(s))
}

func TestIsRestOperation(t *testing.T) {
	assert.True(t, IsRestOperation(createOper("GET")))
}

func TestGetRestOperationMethod(t *testing.T) {
	assert.Equal(t, "GET", GetRestOperationMethod(createOper("GET")))
}

func TestGetRestOperationPath(t *testing.T) {
	assert.Equal(t, "/api/person", GetRestOperationPath(createOper("DONTCARE")))
}

func TestHasInputGet(t *testing.T) {
	assert.False(t, HasInput(createOper("GET")))
}

func TestHasInputDelete(t *testing.T) {
	assert.False(t, HasInput(createOper("DELETE")))
}

func TestHasInputPost(t *testing.T) {
	assert.True(t, HasInput(createOper("POST")))
}

func TestHasInputPut(t *testing.T) {
	assert.True(t, HasInput(createOper("PUT")))
}

func TestGetInputArgTypeString(t *testing.T) {
	o := model.Operation{
		InputArgs: []model.Field{
			{TypeName: "string"},
		},
	}
	assert.Equal(t, "", GetInputArgType(o))
}

func TestGetInputArgTypePerson(t *testing.T) {
	assert.Equal(t, "Person", GetInputArgType(createOper("DONTCARE")))
}

func TestGetInputArgName(t *testing.T) {
	assert.Equal(t, "person", GetInputArgName(createOper("DONTCARE")))
}

func TestGetInputParamString(t *testing.T) {
	assert.Equal(t, "ctx,uid,person", GetInputParamString(createOper("DONTCARE")))
}

func TestHasOutput(t *testing.T) {
	assert.True(t, HasOutput(createOper("DONTCARE")))
}

func TestGetOutputArgType(t *testing.T) {
	assert.Equal(t, "Person", GetOutputArgType(createOper("DONTCARE")))
}

func TestIsPrimitiveTrue(t *testing.T) {
	f := model.Field{Name: "uid", TypeName: "string"}
	assert.True(t, IsPrimitiveArg(f))
}

func TestIsPrimitiveFalse(t *testing.T) {
	f := model.Field{Name: "person", TypeName: "Person"}
	assert.False(t, IsPrimitiveArg(f))
}

func TestIsNumberTrue(t *testing.T) {
	f := model.Field{Name: "uid", TypeName: "int"}
	assert.True(t, IsNumberArg(f))
}

func TestIsNumberFalse(t *testing.T) {
	f := model.Field{Name: "uid", TypeName: "string"}
	assert.False(t, IsNumberArg(f))
}

func createOper(method string) model.Operation {
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
