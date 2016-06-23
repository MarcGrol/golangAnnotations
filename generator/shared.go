package generator

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/MarcGrol/astTools/model"
	"github.com/MarcGrol/astTools/model/annotation"
	"github.com/MarcGrol/astTools/model/annotation/eventAnno"
	"github.com/MarcGrol/astTools/model/annotation/restAnno"
)

var templates map[string]string = map[string]string{
	"aggregates": aggregateTemplate,
	"wrappers":   wrappersTemplate,
	"handlers":   handlersTemplate,
	"helpers":    helpersTemplate,
}

var funcMap = template.FuncMap{
	"IsEvent":                IsEvent,
	"GetAggregateName":       GetAggregateName,
	"IsRestService":          IsRestService,
	"GetRestServicePath":     GetRestServicePath,
	"IsRestOperation":        IsRestOperation,
	"GetRestOperationPath":   GetRestOperationPath,
	"GetRestOperationMethod": GetRestOperationMethod,
	"HasInput":               HasInput,
	"GetInputArgType":        GetInputArgType,
	"GetInputArgName":        GetInputArgName,
	"GetInputParamString":    GetInputParamString,
	"GetOutputArgType":       GetOutputArgType,
	"HasOutput":              HasOutput,
	"IsPrimitive":            IsPrimitive,
	"IsNumber":               IsNumber,
	"ToFirstUpper":           ToFirstUpper,
}

func IsEvent(s model.Struct) bool {
	annotation, ok := annotation.ResolveAnnotations(s.DocLines)
	if !ok || annotation.Name != "Event" {
		return false
	}
	return ok
}

func GetAggregateName(s model.Struct) string {
	val, ok := annotation.ResolveAnnotations(s.DocLines)
	if ok {
		return val.Attributes[eventAnno.ParamAggregate]
	}
	return ""
}

func IsRestService(s model.Struct) bool {
	annotation, ok := annotation.ResolveAnnotations(s.DocLines)
	if !ok || annotation.Name != "RestService" {
		return false
	}
	return ok
}

func GetRestServicePath(o model.Struct) string {
	val, ok := annotation.ResolveAnnotations(o.DocLines)
	if ok {
		return val.Attributes[restAnno.ParamPath]
	}
	return ""
}

func IsRestOperation(o model.Operation) bool {
	annotation, ok := annotation.ResolveAnnotations(o.DocLines)
	if !ok || annotation.Name != "RestOperation" {
		return false
	}
	return ok
}

func GetRestOperationPath(o model.Operation) string {
	val, ok := annotation.ResolveAnnotations(o.DocLines)
	if ok {
		return val.Attributes[restAnno.ParamPath]
	}
	return ""
}

func GetRestOperationMethod(o model.Operation) string {
	val, ok := annotation.ResolveAnnotations(o.DocLines)
	if ok {
		return val.Attributes[restAnno.ParamMethod]
	}
	return ""
}

func HasInput(o model.Operation) bool {
	if GetRestOperationMethod(o) == "POST" || GetRestOperationMethod(o) == "PUT" {
		return true
	}
	return false
}

func GetInputArgType(o model.Operation) string {
	for _, arg := range o.InputArgs {
		if arg.TypeName != "int" && arg.TypeName != "string" {
			return arg.TypeName
		}
	}
	return ""
}

func GetInputArgName(o model.Operation) string {
	for _, arg := range o.InputArgs {
		if arg.TypeName != "int" && arg.TypeName != "string" {
			return arg.Name
		}
	}
	return ""
}

func GetInputParamString(o model.Operation) string {
	args := []string{}
	for _, arg := range o.InputArgs {
		args = append(args, arg.Name)
	}
	return strings.Join(args, ",")
}

func HasOutput(o model.Operation) bool {
	for _, arg := range o.OutputArgs {
		if arg.TypeName != "error" {
			return true
		}
	}
	return false
}

func GetOutputArgType(o model.Operation) string {
	for _, arg := range o.OutputArgs {
		if arg.TypeName != "error" {
			return arg.TypeName
		}
	}
	return ""
}

func IsPrimitive(f model.Field) bool {
	return f.TypeName == "int" || f.TypeName == "string"
}

func IsNumber(f model.Field) bool {
	return f.TypeName == "int"
}

func ToFirstUpper(in string) string {
	if len(in) == 0 {
		return in
	}
	return strings.ToUpper(fmt.Sprintf("%c", in[0])) + in[1:]
}

func getPackageName(structs []model.Struct) (string, error) {
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

func determineTargetPath(inputDir string, packageName string) (string, error) {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		return "", fmt.Errorf("GOPATH not set")
	}
	//log.Printf("GOPATH:%s", goPath)

	workDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("Error getting working dir:%s", err)
	}
	//log.Printf("work-dir:%s", workDir)

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

func generateFileFromTemplate(data interface{}, templateName string, targetFileName string) error {
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
	t, err = t.Parse(templates[templateName])
	if err != nil {
		return err
	}

	defer w.Close()
	if err := t.Execute(w, data); err != nil {
		return err
	}
	return nil
}
