package generator

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/MarcGrol/astTools/model"
)

func GenerateForStruct(myStruct model.Struct) error {
	err := generateEnvelope(myStruct, ".")
	if err != nil {
		return err
	}
	err = generateWrapperForStruct(myStruct, ".")
	if err != nil {
		return err
	}
	return nil
}

func generateEnvelope(str model.Struct, templateDir string) error {
	targetDir := str.PackageName
	dir, _ := path.Split(str.PackageName)
	if dir == "" {
		targetDir = "."
	} else {
		str.PackageName = path.Dir(str.PackageName)
	}
	src := fmt.Sprintf("%s/envelope.go.tmpl", templateDir)
	target := fmt.Sprintf("%s/envelope.go", targetDir)

	err := generateFileFromTemplate(str, src, target)
	if err != nil {
		log.Fatalf("Error generating events (%s)", err)
		return err
	}
	return nil
}

func generateWrapperForStruct(str model.Struct, templateDir string) error {
	if str.IsEvent() {
		targetDir := str.PackageName
		dir, _ := path.Split(str.PackageName)
		if dir == "" {
			targetDir = "."
		} else {
			str.PackageName = path.Dir(str.PackageName)
		}
		src := fmt.Sprintf("%s/wrapper.go.tmpl", templateDir)
		target := fmt.Sprintf("%s/%sWrapperAgain.go", targetDir, str.Name)

		err := generateFileFromTemplate(str, src, target)
		if err != nil {
			log.Fatalf("Error generating events (%s)", err)
			return err
		}
	}
	return nil
}

func generateFileFromTemplate(data interface{}, templateFileName string, targetFileName string) error {
	log.Printf("Using template %s to generate target %s\n", templateFileName, targetFileName)
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(targetFileName), 0777)
	if err != nil {
		return err
	}
	w, err := os.Create(targetFileName)
	if err != nil {
		return err
	}

	defer w.Close()
	if err := t.Execute(w, data); err != nil {
		return err
	}
	return nil
}
