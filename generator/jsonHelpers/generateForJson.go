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

type jsonContext struct {
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
			err := generationUtil.GenerateFileFromTemplateFile(data, packageName, "json-enums", "generator/jsonHelpers/enum.go.tmpl", customTemplateFuncs, target)
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

func IsJsonEnumTolerant(e model.Enum) bool {
	if ann, ok := annotation.ResolveAnnotationByName(e.DocLines, jsonAnnotation.TypeEnum); ok {
		return ann.Attributes[jsonAnnotation.ParamTolerant] == "true"
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

func GetJsonEnumDefault(e model.Enum) string {
	if ann, ok := annotation.ResolveAnnotationByName(e.DocLines, jsonAnnotation.TypeEnum); ok {
		return ann.Attributes[jsonAnnotation.ParamDefault]
	}
	return ""
}

func hasDefaultValue(e model.Enum) bool {
	return GetJsonEnumDefault(e) != ""
}

func getDefaultValue(e model.Enum) string {
	return GetJsonEnumBase(e) + GetJsonEnumDefault(e)
}

func hasAlternativeName(e model.Enum) bool {
	return HasJsonEnumBase(e) && IsJsonEnumTolerant(e)
}

func getAlternativeName(e model.Enum, lit model.EnumLiteral) string {
	if IsJsonEnumStripped(e) {
		return lowerInitial(lit.Name)
	} else {
		base := GetJsonEnumBase(e)
		return lowerInitial(strings.TrimPrefix(lit.Name, base))
	}
}

func getPreferredName(e model.Enum, lit model.EnumLiteral) string {
	if IsJsonEnumStripped(e) {
		base := GetJsonEnumBase(e)
		return lowerInitial(strings.TrimPrefix(lit.Name, base))
	} else {
		return lowerInitial(lit.Name)
	}
}

func lowerInitial(s string) string {
	a := []rune(s)
	a[0] = unicode.ToLower(a[0])
	return string(a)
}

func IsJsonStruct(s model.Struct) bool {
	_, ok := annotation.ResolveAnnotationByName(s.DocLines, jsonAnnotation.TypeStruct)
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
