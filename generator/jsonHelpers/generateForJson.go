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

type Generator struct {
}

func NewGenerator() generationUtil.Generator {
	return &Generator{}
}

func (eg *Generator) GetAnnotations() []annotation.AnnotationDescriptor {
	return jsonAnnotation.Get()
}

func (eg *Generator) Generate(inputDir string, parsedSource model.ParsedSources) error {
	return generate(inputDir, parsedSource.Enums, parsedSource.Structs)
}

type jsonContext struct {
	PackageName string
	Enums       []model.Enum
	Structs     []model.Struct
}

func generate(inputDir string, enums []model.Enum, structs []model.Struct) error {

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
		if IsJSONEnum(anEnum) {
			jsonEnums = append(jsonEnums, anEnum)
		}
	}
	jsonStructs := make([]model.Struct, 0, len(structs))
	for _, aStruct := range structs {
		if IsJSONStruct(aStruct) {
			jsonStructs = append(jsonStructs, aStruct)
		}
	}
	if len(jsonEnums) == 0 && len(jsonStructs) == 0 {
		return nil
	}

	err = doGenerate(packageName, jsonEnums, jsonStructs, targetDir)
	if err != nil {
		return err
	}

	return nil
}

func doGenerate(packageName string, jsonEnums []model.Enum, jsonStructs []model.Struct, targetDir string) error {
	filenameMap := getFilenamesWithTypeNames(jsonEnums, jsonStructs)

	for fn := range filenameMap {
		targetFilename := strings.Replace(fn, ".", "_json.", 1)
		target := fmt.Sprintf("%s/$%s", targetDir, targetFilename)

		data := jsonContext{
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
			err := generationUtil.GenerateFileFromTemplate(data, packageName, "json-enums", jsonHelpersTemplate, customTemplateFuncs, target)
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
	"HasAlternativeName": hasAlternativeName,
	"GetAlternativeName": getAlternativeName,
	"GetPreferredName":   getPreferredName,
	"HasDefaultValue":    hasDefaultValue,
	"GetDefaultValue":    getDefaultValue,
	"HasSlices":          hasSlices,
}

func IsJSONEnum(e model.Enum) bool {
	annotations := annotation.NewRegistry(jsonAnnotation.Get())
	_, ok := annotations.ResolveAnnotationByName(e.DocLines, jsonAnnotation.TypeEnum)
	return ok
}

func IsJSONEnumStripped(e model.Enum) bool {
	annotations := annotation.NewRegistry(jsonAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(e.DocLines, jsonAnnotation.TypeEnum); ok {
		return ann.Attributes[jsonAnnotation.ParamStripped] == "true"
	}
	return false
}

func IsJSONEnumTolerant(e model.Enum) bool {
	annotations := annotation.NewRegistry(jsonAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(e.DocLines, jsonAnnotation.TypeEnum); ok {
		return ann.Attributes[jsonAnnotation.ParamTolerant] == "true"
	}
	return false
}

func GetJSONEnumBase(e model.Enum) string {
	annotations := annotation.NewRegistry(jsonAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(e.DocLines, jsonAnnotation.TypeEnum); ok {
		return ann.Attributes[jsonAnnotation.ParamBase]
	}
	return ""
}

func HasJSONEnumBase(e model.Enum) bool {
	return GetJSONEnumBase(e) != ""
}

func GetJSONEnumDefault(e model.Enum) string {
	annotations := annotation.NewRegistry(jsonAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(e.DocLines, jsonAnnotation.TypeEnum); ok {
		return ann.Attributes[jsonAnnotation.ParamDefault]
	}
	return ""
}

func hasDefaultValue(e model.Enum) bool {
	return GetJSONEnumDefault(e) != ""
}

func getDefaultValue(e model.Enum) string {
	return GetJSONEnumBase(e) + GetJSONEnumDefault(e)
}

func hasAlternativeName(e model.Enum) bool {
	return HasJSONEnumBase(e) && IsJSONEnumTolerant(e)
}

func getAlternativeName(e model.Enum, lit model.EnumLiteral) string {
	if IsJSONEnumStripped(e) {
		return lowerInitial(lit.Name)
	}
	base := GetJSONEnumBase(e)
	return lowerInitial(strings.TrimPrefix(lit.Name, base))
}

func getPreferredName(e model.Enum, lit model.EnumLiteral) string {
	if IsJSONEnumStripped(e) {
		base := GetJSONEnumBase(e)
		return lowerInitial(strings.TrimPrefix(lit.Name, base))
	}
	return lowerInitial(lit.Name)
}

func lowerInitial(s string) string {
	a := []rune(s)
	a[0] = unicode.ToLower(a[0])
	return string(a)
}

func IsJSONStruct(s model.Struct) bool {
	annotations := annotation.NewRegistry(jsonAnnotation.Get())
	_, ok := annotations.ResolveAnnotationByName(s.DocLines, jsonAnnotation.TypeStruct)
	return ok
}

func hasSlices(s model.Struct) bool {
	for _, f := range s.Fields {
		if f.IsSlice {
			return true
		}
	}
	return false
}
