package generationUtil

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/MarcGrol/golangAnnotations/model"
)

func GetPackageName(structs []model.Struct) (string, error) {
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

func GenerateFileFromTemplate(data interface{}, templateName string, templateString string, funcMap template.FuncMap, targetFileName string) error {
	log.Printf("Using template '%s' to generate target %s\n", templateName, targetFileName)

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
