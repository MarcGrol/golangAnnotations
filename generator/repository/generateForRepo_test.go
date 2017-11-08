package repository

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/MarcGrol/golangAnnotations/model"
	"github.com/stretchr/testify/assert"
)

func cleanup() {
	os.Remove("./testData/$userRepo.go")
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
	_, err = os.Stat("./testData/$userRepo.go")
	assert.NoError(t, err)

	// check that generate code has 4 helper functions for MyStruct
	data, err := ioutil.ReadFile("./testData/$userRepo.go")
	assert.NoError(t, err)
	assert.Contains(t, string(data), `func DefaultFindEndUserOnUID(c context.Context, credentials rest.Credentials, endUserUID string) (*model.EndUser, error) {
`)
}
