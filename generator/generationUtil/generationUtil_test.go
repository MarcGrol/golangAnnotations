package generationUtil

import (
	"io/ioutil"
	"os"
	"testing"
	"text/template"

	"github.com/MarcGrol/golangAnnotations/model"
	"github.com/stretchr/testify/assert"
)

func TestGetPackageName(t *testing.T) {
	s := []model.Struct{
		{PackageName: "mypack"},
	}
	packName, err := GetPackageNameForStructs(s)
	assert.Equal(t, "mypack", packName)
	assert.Nil(t, err)
}

func TestGetPackageNameMultiplePackages(t *testing.T) {
	s := []model.Struct{
		{PackageName: "mypack"},
		{PackageName: "otheePack"},
	}
	_, err := GetPackageNameForStructs(s)
	assert.Error(t, err)
}

func TestGetPackageNameNoStructs(t *testing.T) {
	s := []model.Struct{}
	name, err := GetPackageNameForStructs(s)
	assert.Empty(t, name)
	assert.NoError(t, err)

}

func TestDetermineTargetPathEmptyInput(t *testing.T) {
	inputDir := ""
	packageName := ""
	_, err := DetermineTargetPath(inputDir, packageName)
	assert.Error(t, err)
	assert.Equal(t, "Input params not set", err.Error())
}

func TestDetermineTargetCurrent(t *testing.T) {
	inputDir := "."
	packageName := "generationUtil"
	dir, _ := DetermineTargetPath(inputDir, packageName)
	assert.Equal(t, ".", dir)
}

func TestDetermineTargetSubdir(t *testing.T) {
	inputDir := "a/b"
	packageName := "generationUtil"
	dir, _ := DetermineTargetPath(inputDir, packageName)
	assert.Equal(t, "a/b/generationUtil", dir)
}

func CommentedPackageName(s model.Struct) string {
	return "// commented " + s.PackageName
}

func TestGenerateFileFromTemplate(t *testing.T) {
	var fm = template.FuncMap{
		"CommentedPackageName": CommentedPackageName,
	}

	err := Generate(Info{
		Src:            "testsrc",
		TargetFilename: "test/doit.txt",
		TemplateName:   "testtemplate",
		TemplateString: "{{.PackageName}}\n{{CommentedPackageName .}}",
		FuncMap:        fm,
		Data:           model.Struct{PackageName: "testit"},
	})
	assert.Nil(t, err)

	data, err := ioutil.ReadFile("test/doit.txt")
	assert.NoError(t, err)

	assert.Equal(t, "testit\n// commented testit", string(data))

	os.Remove("./test/doit.txt")
	os.Remove("./test")
}
