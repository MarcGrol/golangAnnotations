package jsonHelpers

import (
	"fmt"
	"log"
	"text/template"

	"github.com/MarcGrol/golangAnnotations/annotation"
	"github.com/MarcGrol/golangAnnotations/generator/generationUtil"
	"github.com/MarcGrol/golangAnnotations/generator/jsonHelpers/jsonAnnotation"
	"github.com/MarcGrol/golangAnnotations/model"
)

type Enums struct {
	PackageName string
	Enums       []model.Enum
	Structs     []model.Struct
}

func Generate(inputDir string, parsedSource model.ParsedSources) error {
	return generate(inputDir, parsedSource.Enums, parsedSource.Structs)
}

func generate(inputDir string, enums []model.Enum, structs []model.Struct) error {
	jsonAnnotation.Register()

	packageName, err := generationUtil.GetPackageNameForEnums(enums)
	if err != nil {
		return err
	}
	targetDir, err := generationUtil.DetermineTargetPath(inputDir, packageName)
	if err != nil {
		return err
	}
	target := fmt.Sprintf("%s/jsonHelpers.go", targetDir)

	data := Enums{
		PackageName: packageName,
		Enums:       enums,
		Structs:     structs,
	}
	err = generationUtil.GenerateFileFromTemplate(data, packageName, "enums", enumTemplate, customTemplateFuncs, target)
	if err != nil {
		log.Fatalf("Error generating wrappers for enums (%s)", err)
		return err
	}

	return nil
}

var customTemplateFuncs = template.FuncMap{
	"IsJsonEnum":   IsJsonEnum,
	"IsJsonStruct": IsJsonStruct,
	"HasSlices":    HasSlices,
}

func IsJsonEnum(e model.Enum) bool {
	annotation, ok := annotation.ResolveAnnotations(e.DocLines)
	if !ok || annotation.Name != "JsonEnum" {
		return false
	}
	return ok
}

func IsJsonStruct(s model.Struct) bool {
	annotation, ok := annotation.ResolveAnnotations(s.DocLines)
	if !ok || annotation.Name != "JsonStruct" {
		return false
	}
	return ok
}

func HasSlices(s model.Struct) bool {
	for _, f := range s.Fields {
		if f.IsSlice {
			return true
		}
	}
	return false
}

var enumTemplate string = `
// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

{{range .Enums}}

{{if IsJsonEnum . }}

// Helpers for json-enum {{.Name}}

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

// MarshalJSON caters for readable enums with a proper default value
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

// UnmarshalJSON caters for readable enums with a proper default value
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


{{range .Structs}}

{{if IsJsonStruct . }}

// Helpers for json-struct {{.Name}}

{{if HasSlices . }}

// MarshalJSON prevents nil slices in json
func (data {{.Name}}) MarshalJSON() ([]byte, error) {
	type alias {{.Name}}
	var raw = alias(data)

	{{range .Fields}}
		{{if .IsSlice}}
		if raw.{{.Name}} == nil {
			raw.{{.Name}} = []{{if .IsPointer}}*{{end}}{{.TypeName}}{}
		}
		{{end}}
	{{end}}

	return json.Marshal(raw)
}

// UnmarshalJSON prevents nil slices from json
func (data *{{.Name}}) UnmarshalJSON(b []byte) error {
	type alias {{.Name}}
	var raw alias
	err := json.Unmarshal(b, &raw)

	{{range .Fields}}
		{{if .IsSlice}}
		if raw.{{.Name}} == nil {
			raw.{{.Name}} = []{{if .IsPointer}}*{{end}}{{.TypeName}}{}
		}
		{{end}}
	{{end}}

	*data = {{.Name}}(raw)

	return err
}

{{end}}

{{end}}
{{end}}
`
