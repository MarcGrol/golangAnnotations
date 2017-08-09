package jsonHelpers

import (
	"fmt"
	"log"
	"strings"
	"text/template"
	"unicode"

	"github.com/MarcGrol/golangAnnotations/annotation"
	"github.com/MarcGrol/golangAnnotations/generator/generationUtil"
	"github.com/MarcGrol/golangAnnotations/generator/jsonHelpers/jsonAnnotation"
	"github.com/MarcGrol/golangAnnotations/model"
)

type JsonContext struct {
	PackageName string
	Enums       []model.Enum
	Structs     []model.Struct
}

func Generate(inputDir string, parsedSource model.ParsedSources) error {
	return generate(inputDir, parsedSource.Enums, parsedSource.Structs)
}

func generate(inputDir string, enums []model.Enum, structs []model.Struct) error {
	jsonAnnotation.Register()

	packageName, err := generationUtil.GetPackageNameForEnumsOrStructs(enums, structs)
	if err != nil {
		return err
	}
	targetDir, err := generationUtil.DetermineTargetPath(inputDir, packageName)
	if err != nil {
		return err
	}

	jsonEnums := make([]model.Enum, 0, len(enums))
	for _, anEnum := range enums {
		if IsJsonEnum(anEnum) {
			jsonEnums = append(jsonEnums, anEnum)
		}
	}
	jsonStructs := make([]model.Struct, 0, len(structs))
	for _, aStruct := range structs {
		if IsJsonStruct(aStruct) {
			jsonStructs = append(jsonStructs, aStruct)
		}
	}
	if len(jsonEnums) == 0 && len(jsonStructs) == 0 {
		return nil
	}

	filenameMap := getFilenamesWithTypeNames(jsonEnums, jsonStructs)

	for fn := range filenameMap {
		targetFilename := strings.Replace(fn, ".", "_json.", 1)
		target := fmt.Sprintf("%s/$%s", targetDir, targetFilename)

		data := JsonContext{
			PackageName: packageName,
		}

		// find al enums belonging to this file
		for _, e := range jsonEnums {
			if e.Filename == fn {
				data.Enums = append(data.Enums, e)
			}
		}
		for _, s := range jsonStructs {
			if s.Filename == fn {
				data.Structs = append(data.Structs, s)
			}
		}

		if len(data.Enums) > 0 || len(data.Structs) > 0 {
			err = generationUtil.GenerateFileFromTemplate(data, packageName, "enums", enumTemplate, customTemplateFuncs, target)
			if err != nil {
				log.Fatalf("Error generating wrappers for enums (%s)", err)
				return err
			}
		}
	}

	return nil
}

func getFilenamesWithTypeNames(jsonEnums []model.Enum, jsonStructs []model.Struct) map[string][]string {
	// group enum and structs by filename
	filenameMap := map[string][]string{}

	// get all  enum-names belonging to file
	for _, e := range jsonEnums {
		typeNames := filenameMap[e.Filename]
		typeNames = append(typeNames, e.Name)
		filenameMap[e.Filename] = typeNames
	}

	// get all  struct-names belonging to file
	for _, s := range jsonStructs {
		typeNames := filenameMap[s.Filename]
		typeNames = append(typeNames, s.Name)
		filenameMap[s.Filename] = typeNames
	}

	return filenameMap
}

var customTemplateFuncs = template.FuncMap{
	"HasJsonEnumBase": HasJsonEnumBase,
	"UnstrippedName":  UnstrippedName,
	"StrippedName":    StrippedName,
	"PreferredName":   PreferredName,
	"HasSlices":       HasSlices,
}

func IsJsonEnum(e model.Enum) bool {
	_, ok := annotation.ResolveAnnotationByName(e.DocLines, jsonAnnotation.TypeEnum)
	return ok
}

func IsJsonEnumStripped(e model.Enum) bool {
	if ann, ok := annotation.ResolveAnnotationByName(e.DocLines, jsonAnnotation.TypeEnum); ok {
		return ann.Attributes[jsonAnnotation.ParamStripped] == "true"
	}
	return false
}

func GetJsonEnumBase(e model.Enum) string {
	if ann, ok := annotation.ResolveAnnotationByName(e.DocLines, jsonAnnotation.TypeEnum); ok {
		return ann.Attributes[jsonAnnotation.ParamBase]
	}
	return ""
}

func HasJsonEnumBase(e model.Enum) bool {
	return GetJsonEnumBase(e) != ""
}

func IsJsonStruct(s model.Struct) bool {
	_, ok := annotation.ResolveAnnotationByName(s.DocLines, jsonAnnotation.TypeStruct)
	return ok
}

func UnstrippedName(lit model.EnumLiteral) string {
	return lowerInitial(lit.Name)
}

func StrippedName(e model.Enum, lit model.EnumLiteral) string {
	base := GetJsonEnumBase(e)
	return lowerInitial(strings.TrimPrefix(lit.Name, base))
}

func PreferredName(e model.Enum, lit model.EnumLiteral) string {
	if IsJsonEnumStripped(e) {
		return StrippedName(e, lit)
	} else {
		return UnstrippedName(lit)
	}
}

func lowerInitial(s string) string {
	a := []rune(s)
	a[0] = unicode.ToLower(a[0])
	return string(a)
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

import "encoding/json"

{{range .Enums}}
{{$enum := .}}

// Helpers for json-enum {{.Name}}

var (
	_{{.Name}}NameToValue = map[string]{{.Name}}{
		{{range .EnumLiterals}}"{{UnstrippedName .}}":{{.Name}},
		{{end}}
		{{if HasJsonEnumBase $enum}}{{range .EnumLiterals}}"{{StrippedName $enum .}}":{{.Name}},
		{{end}}{{end}}
	}
	_{{.Name}}ValueToName = map[{{.Name}}]string{
		{{range .EnumLiterals }}{{.Name}}:"{{PreferredName $enum .}}",
		{{end}}
	}
)

func {{.Name}}ByName(name string, unknown {{.Name}}) {{.Name}} {
	t, ok := _{{.Name}}NameToValue[name]
	if !ok {
		return unknown
	}
	return t
}

func (t {{.Name}}) String() string {
	v, _ := _{{.Name}}ValueToName[t]
	return v
}

// MarshalJSON caters for readable enums with a proper default value
func (r {{.Name}}) MarshalJSON() ([]byte, error) {
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

{{range .Structs}}

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
`
