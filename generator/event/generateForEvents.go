package event

import (
	"fmt"
	"log"
	"text/template"

	"github.com/MarcGrol/golangAnnotations/annotation"
	"github.com/MarcGrol/golangAnnotations/generator/event/eventAnnotation"
	"github.com/MarcGrol/golangAnnotations/generator/generationUtil"
	"github.com/MarcGrol/golangAnnotations/model"
)

type AggregateMap struct {
	PackageName  string
	AggregateMap map[string]map[string]string
}

type Structs struct {
	PackageName string
	Structs     []model.Struct
}

func Generate(inputDir string, parsedSource model.ParsedSources) error {
	return generate(inputDir, parsedSource.Structs)
}

func generate(inputDir string, structs []model.Struct) error {
	eventAnnotation.Register()

	packageName, err := generationUtil.GetPackageNameForStructs(structs)
	if err != nil {
		return err
	}
	aggregates := make(map[string]map[string]string)
	eventCount := 0
	for _, s := range structs {
		if IsEvent(s) {
			events, ok := aggregates[GetAggregateName(s)]
			if !ok {
				events = make(map[string]string)
			}
			events[s.Name] = s.Name
			aggregates[GetAggregateName(s)] = events
			eventCount++
		}
	}

	if eventCount > 0 {
		targetDir, err := generationUtil.DetermineTargetPath(inputDir, packageName)
		if err != nil {
			return err
		}
		{
			target := fmt.Sprintf("%s/$aggregates.go", targetDir)

			data := AggregateMap{
				PackageName:  packageName,
				AggregateMap: aggregates,
			}

			err = generationUtil.GenerateFileFromTemplateFile(data, packageName, "aggregates", "generator/event/aggregate.go.tmpl", customTemplateFuncs, target)
			if err != nil {
				log.Fatalf("Error generating aggregates (%s)", err)
				return err
			}
		}
		{
			target := fmt.Sprintf("%s/$wrappers.go", targetDir)

			data := Structs{
				PackageName: packageName,
				Structs:     structs,
			}
			err = generationUtil.GenerateFileFromTemplateFile(data, packageName, "wrappers", "generator/event/wrappers.go.tmpl", customTemplateFuncs, target)
			if err != nil {
				log.Fatalf("Error generating wrappers for structs (%s)", err)
				return err
			}
		}
		{
			target := fmt.Sprintf("%s/../store/$%sEventStore.go", targetDir, packageName)

			data := Structs{
				PackageName: packageName,
				Structs:     structs,
			}
			err = generationUtil.GenerateFileFromTemplateFile(data, packageName, "store-events", "generator/event/eventStore.go.tmpl", customTemplateFuncs, target)
			if err != nil {
				log.Fatalf("Error generating store-events for structs (%s)", err)
				return err
			}
		}
		{
			target := fmt.Sprintf("%s/$wrappers_test.go", targetDir)

			data := Structs{
				PackageName: packageName,
				Structs:     structs,
			}
			err = generationUtil.GenerateFileFromTemplateFile(data, packageName, "wrappers-test", "generator/event/wrappers_test.go.tmpl", customTemplateFuncs, target)
			if err != nil {
				log.Fatalf("Error generating wrappers-test for structs (%s)", err)
				return err
			}
		}
	}
	return nil
}

var customTemplateFuncs = template.FuncMap{
	"IsEvent":          IsEvent,
	"IsRootEvent":      IsRootEvent,
	"IsPersistent":     IsPersistent,
	"IsTransient":      IsTransient,
	"GetAggregateName": GetAggregateName,
	"HasValueForField": HasValueForField,
	"ValueForField":    ValueForField,
}

func IsEvent(s model.Struct) bool {
	_, ok := annotation.ResolveAnnotationByName(s.DocLines, eventAnnotation.TypeEvent)
	return ok
}

func GetAggregateName(s model.Struct) string {
	if ann, ok := annotation.ResolveAnnotationByName(s.DocLines, eventAnnotation.TypeEvent); ok {
		return ann.Attributes[eventAnnotation.ParamAggregate]
	}
	return ""
}

func IsRootEvent(s model.Struct) bool {
	if ann, ok := annotation.ResolveAnnotationByName(s.DocLines, eventAnnotation.TypeEvent); ok {
		return ann.Attributes[eventAnnotation.ParamIsRootEvent] == "true"
	}
	return false
}

func IsPersistent(s model.Struct) bool {
	return !IsTransient(s)
}

func IsTransient(s model.Struct) bool {
	if ann, ok := annotation.ResolveAnnotationByName(s.DocLines, eventAnnotation.TypeEvent); ok {
		return ann.Attributes[eventAnnotation.ParamIsTransient] == "true"
	}
	return false
}

func HasValueForField(field model.Field) bool {
	if field.TypeName == "int" || field.TypeName == "string" || field.TypeName == "bool" {
		return true
	}
	return false
}

func ValueForField(field model.Field) string {
	if field.TypeName == "int" {
		if field.IsSlice {
			return "[]int{1,2}"
		} else {
			return "42"
		}
	} else if field.TypeName == "string" {
		if field.IsSlice {
			return "[]string{" + fmt.Sprintf("\"Example1%s\"", field.Name) + "," +
				fmt.Sprintf("\"Example1%s\"", field.Name) + "}"
		} else {
			return fmt.Sprintf("\"Example3%s\"", field.Name)
		}
	} else if field.TypeName == "bool" {
		return "true"
	}
	return ""
}
