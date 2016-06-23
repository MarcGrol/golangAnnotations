package generator

import (
	"os"
	"testing"

	"io/ioutil"

	"github.com/MarcGrol/astTools/model"
	"github.com/stretchr/testify/assert"
)

func TestGenerateForWeb(t *testing.T) {
	os.Remove("./testData/httpMyService.go")
	os.Remove("./testData/httpMyServiceHelpers_test.go")

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
			DocLines:      []string{"// @RestOperation(path = \"/person\", method = \"GET\")"},
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

	err := GenerateForWeb("testData", s)
	assert.Nil(t, err)

	// check that generated files exisst
	_, err = os.Stat("./testData/httpMyService.go")
	assert.NoError(t, err)

	// check that generate code has 4 helper functions for MyStruct
	data, err := ioutil.ReadFile("./testData/httpMyService.go")
	assert.NoError(t, err)
	assert.Contains(t, string(data), "func (ts *MyService) HttpHandler() http.Handler {")
	assert.Contains(t, string(data), "func doit( service *MyService ) http.HandlerFunc {")

	// check that generated files exisst
	_, err = os.Stat("./testData/httpMyService.go")
	assert.NoError(t, err)

	// check that generate code has 4 helper functions for MyStruct
	data, err = ioutil.ReadFile("./testData/httpMyServiceHelpers_test.go")
	assert.NoError(t, err)
	assert.Contains(t, string(data), "func doitTestHelper")

	os.Remove("./testData/httpMyService.go")
	os.Remove("./testData/httpMyServiceHelpers_test.go")

}
