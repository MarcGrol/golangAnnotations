package repository

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/f0rt/golangAnnotations/generator/generationUtil"
	"github.com/f0rt/golangAnnotations/model"
	"github.com/stretchr/testify/assert"
)

func cleanup() {
	os.Remove(generationUtil.Prefixed("./testData/ast.json"))
	os.Remove(generationUtil.Prefixed("./testData/userRepo.go"))
}

func TestGenerateForRepo(t *testing.T) {
	cleanup()
	defer cleanup()

	s := []model.Struct{
		{
			DocLines:    []string{`// @Repository( aggregate = "User", model="EndUser", package="testEvents", methods="find" )`},
			PackageName: "testData",
			Name:        "UserRepo",
		},
	}

	err := NewGenerator().Generate("testData", model.ParsedSources{Structs: s})
	assert.Nil(t, err)

	// check that generated files exisst
	_, err = os.Stat(generationUtil.Prefixed("./testData/userRepo.go"))
	assert.NoError(t, err)

	// check that generate code has 4 helper functions for MyStruct
	data, err := ioutil.ReadFile(generationUtil.Prefixed("./testData/userRepo.go"))
	assert.NoError(t, err)
	assert.Contains(t, string(data), `func DefaultFindEndUserOnUID(c context.Context, rc request.Context, tx *datastore.Transaction, endUserUID string) (*endUserModel.EndUser, error) {
`)
}
