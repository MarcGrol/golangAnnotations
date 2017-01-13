package enum

import (
	"fmt"
	"log"
	"text/template"

	"github.com/MarcGrol/golangAnnotations/annotation"
	"github.com/MarcGrol/golangAnnotations/generator/enum/enumAnnotation"
	"github.com/MarcGrol/golangAnnotations/generator/generationUtil"
	"github.com/MarcGrol/golangAnnotations/model"
)

type Enums struct {
	PackageName string
	Enums       []model.Enum
}

func Generate(inputDir string, parsedSource model.ParsedSources) error {
	return generate(inputDir, parsedSource.Enums)
}

func generate(inputDir string, enums []model.Enum) error {
	enumAnnotation.Register()

	packageName, err := generationUtil.GetPackageNameForEnums(enums)
	if err != nil {
		return err
	}
	targetDir, err := generationUtil.DetermineTargetPath(inputDir, packageName)
	if err != nil {
		return err
	}
	target := fmt.Sprintf("%s/enums.go", targetDir)

	data := Enums{
		PackageName: packageName,
		Enums:       enums,
	}
	err = generationUtil.GenerateFileFromTemplate(data, packageName, "enums", enumTemplate, customTemplateFuncs, target)
	if err != nil {
		log.Fatalf("Error generating wrappers for enums (%s)", err)
		return err
	}

	return nil
}

var customTemplateFuncs = template.FuncMap{
	"IsEnum": IsEnum,
}

func IsEnum(e model.Enum) bool {
	annotation, ok := annotation.ResolveAnnotations(e.DocLines)
	if !ok || annotation.Name != "Enum" {
		return false
	}
	return ok
}

var enumTemplate string = `
// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}


{{range .Enums}}

{{if IsEnum . }}

var (
	_{{.Name}}NameToValue = map[string]{{.Name}}{
		{{range .EnumLiterals}}
		"{{.Name}}":{{.Name}},
		{{end}}
	}

	_{{.Name}}ValueToName = map[{{.Name}}]string{
		{{range .EnumLiterals }}
		{{.Name}}:"{{.Name}}",
		{{end}}
	}
)

func init() {
	var v {{.Name}}
	if _, ok := interface{}(v).(fmt.Stringer); ok {
		_{{.Name}}NameToValue = map[string]{{.Name}}{
			{{range .EnumLiterals }}
			interface{}({{.Name}}).(fmt.Stringer).String():  {{.Name}},
			{{end}}
		}
	}
}

// MarshalJSON is generated so color satisfies json.Marshaler.
func (r {{.Name}}) MarshalJSON() ([]byte, error) {
	if s, ok := interface{}(r).(fmt.Stringer); ok {
		return json.Marshal(s.String())
	}
	s, ok := _{{.Name}}ValueToName[r]
	if !ok {
		return nil, fmt.Errorf("invalid {{.Name}}: %d", r)
	}
	return json.Marshal(s)
}

// UnmarshalJSON is generated so color satisfies json.Unmarshaler.
func (r *{{.Name}}) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("{{.Name}} should be a string, got %s", data)
	}
	v, ok := _{{.Name}}NameToValue[s]
	if !ok {
		return fmt.Errorf("invalid {{.Name}} %q", s)
	}
	*r = v
	return nil
}


{{end}}
{{end}}
`
