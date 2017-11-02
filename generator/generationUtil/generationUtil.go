package generationUtil

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"io/ioutil"

	"github.com/MarcGrol/golangAnnotations/model"
)

type Generator interface {
	Generate(inputDir string, parsedSources model.ParsedSources) error
}

func GetPackageNameForStructs(structs []model.Struct) (string, error) {
	if len(structs) == 0 {
		return "", fmt.Errorf("Need at least one struct to determine package-name")
	}
	packageName := structs[0].PackageName
	for _, s := range structs {
		if s.PackageName != packageName {
			return "", fmt.Errorf("List of structs has multiple package-names")
		}
	}
	return packageName, nil
}

func getPackageNameForEnums(enums []model.Enum) (string, error) {
	if len(enums) == 0 {
		return "", fmt.Errorf("Need at least one enum to determine package-name")
	}
	packageName := enums[0].PackageName
	for _, s := range enums {
		if s.PackageName != packageName {
			return "", fmt.Errorf("List of enums has multiple package-names")
		}
	}
	return packageName, nil
}

func GetPackageNameForEnumsOrStructs(enums []model.Enum, structs []model.Struct) (string, error) {
	if len(enums) == 0 && len(structs) == 0 {
		return "", fmt.Errorf("Need at least one enum or struct to determine package-name")
	}
	var packageNameEnums, packageNameStructs string
	var err error
	if len(enums) > 0 {
		packageNameEnums, err = getPackageNameForEnums(enums)
		if err != nil {
			return "", err
		}
	}
	if len(structs) > 0 {
		packageNameStructs, err = GetPackageNameForStructs(structs)
		if err != nil {
			return "", err
		}
	}
	if packageNameEnums == packageNameStructs || packageNameStructs == "" {
		return packageNameEnums, nil
	}
	if packageNameEnums == "" {
		return packageNameStructs, nil
	}
	return "", fmt.Errorf("List of enums and structs has multiple package-names")
}

func DetermineTargetPath(inputDir string, packageName string) (string, error) {
	if inputDir == "" || packageName == "" {
		return "", fmt.Errorf("Input params not set")
	}

	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		return "", fmt.Errorf("GOPATH not set")
	}

	workDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("Error getting working dir:%s", err)
	}

	if !strings.Contains(workDir, goPath) {
		return "", fmt.Errorf("Code %s lives outside GOPATH:%s", workDir, goPath)
	}

	baseDir := path.Base(inputDir)
	if baseDir == "." || baseDir == packageName {
		return inputDir, nil
	} else {
		return fmt.Sprintf("%s/%s", inputDir, packageName), nil
	}
}

func GenerateFileFromTemplateFile(data interface{}, srcName string, templateName string, templateFilePath string, funcMap template.FuncMap, targetFileName string) error {
	goPath := os.Getenv("GOPATH")
	appPath := goPath + "/src/github.com/MarcGrol/golangAnnotations/"

	cwd, _ := os.Getwd()
	templateString, err := ioutil.ReadFile(appPath + templateFilePath)
	if err != nil {
		log.Printf("Error loading template file from dir %s: %s", cwd, err)
		return err
	}
	return generateFileFromTemplate(data, srcName, templateName, string(templateString), funcMap, targetFileName)
}

func generateFileFromTemplate(data interface{}, srcName string, templateName string, templateString string, funcMap template.FuncMap, targetFileName string) error {
	fmt.Fprintf(os.Stderr, "%s: Generated go file '%s' based on source '%s'\n", "golangAnnotations", targetFileName, srcName)

	err := os.MkdirAll(filepath.Dir(targetFileName), 0777)
	if err != nil {
		return err
	}
	w, err := os.Create(targetFileName)
	if err != nil {
		return err
	}

	t := template.New(templateName).Funcs(funcMap)
	t, err = t.Parse(templateString)
	if err != nil {
		return err
	}

	defer w.Close()
	if err := t.Execute(w, data); err != nil {
		return err
	}
	return nil
}
