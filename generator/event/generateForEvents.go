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

type aggregateMap struct {
	PackageName  string
	AggregateMap map[string]map[string]string
}

type structures struct {
	PackageName string
	Structs     []model.Struct
}

type Generator struct {
}

func NewGenerator() generationUtil.Generator {
	return &Generator{}
}

func (eg *Generator) GetAnnotations() []annotation.AnnotationDescriptor {
	return eventAnnotation.Get()
}

var annotations annotation.AnnotationRegister

func registerAnnotations() {
	annotations = annotation.NewRegistry(eventAnnotation.Get())
}

func (eg *Generator) Generate(inputDir string, parsedSource model.ParsedSources) error {
	registerAnnotations()
	return generate(inputDir, parsedSource.Structs)
}

func generate(inputDir string, structs []model.Struct) error {
	packageName, err := generationUtil.GetPackageNameForStructs(structs)
	if err != nil {
		return err
	}
	aggregates := make(map[string]map[string]string)
	eventCount := 0
	for _, s := range structs {
		if isEvent(s) {
			events, ok := aggregates[getAggregateName(s)]
			if !ok {
				events = make(map[string]string)
			}
			events[s.Name] = s.Name
			aggregates[getAggregateName(s)] = events
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

			data := aggregateMap{
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

			data := structures{
				PackageName: packageName,
				Structs:     structs,
			}
			err = generationUtil.GenerateFileFromTemplateFile(data, packageName, "wrappers", "generator/event/wrappers.go.tmpl", customTemplateFuncs, target)
			if err != nil {
				log.Fatalf("Error generating wrappers for structures (%s)", err)
				return err
			}
		}
		{
			target := fmt.Sprintf("%s/../store/$%sEventStore.go", targetDir, packageName)

			data := structures{
				PackageName: packageName,
				Structs:     structs,
			}
			err = generationUtil.GenerateFileFromTemplateFile(data, packageName, "store-events", "generator/event/eventStore.go.tmpl", customTemplateFuncs, target)
			if err != nil {
				log.Fatalf("Error generating store-events for structures (%s)", err)
				return err
			}
		}
		{
			target := fmt.Sprintf("%s/$wrappers_test.go", targetDir)

			data := structures{
				PackageName: packageName,
				Structs:     structs,
			}
			err = generationUtil.GenerateFileFromTemplateFile(data, packageName, "wrappers-test", "generator/event/wrappers_test.go.tmpl", customTemplateFuncs, target)
			if err != nil {
				log.Fatalf("Error generating wrappers-test for structures (%s)", err)
				return err
			}
		}
	}
	return nil
}

var customTemplateFuncs = template.FuncMap{
	"IsEvent":          isEvent,
	"IsRootEvent":      isRootEvent,
	"IsPersistent":     isPersistent,
	"IsTransient":      isTransient,
	"GetAggregateName": getAggregateName,
	"HasValueForField": hasValueForField,
	"ValueForField":    valueForField,
}

func isEvent(s model.Struct) bool {
	_, ok := annotations.ResolveAnnotationByName(s.DocLines, eventAnnotation.TypeEvent)
	return ok
}

func getAggregateName(s model.Struct) string {
	if ann, ok := annotations.ResolveAnnotationByName(s.DocLines, eventAnnotation.TypeEvent); ok {
		return ann.Attributes[eventAnnotation.ParamAggregate]
	}
	return ""
}

func isRootEvent(s model.Struct) bool {
	if ann, ok := annotations.ResolveAnnotationByName(s.DocLines, eventAnnotation.TypeEvent); ok {
		return ann.Attributes[eventAnnotation.ParamIsRootEvent] == "true"
	}
	return false
}

func isPersistent(s model.Struct) bool {
	return !isTransient(s)
}

func isTransient(s model.Struct) bool {
	if ann, ok := annotations.ResolveAnnotationByName(s.DocLines, eventAnnotation.TypeEvent); ok {
		return ann.Attributes[eventAnnotation.ParamIsTransient] == "true"
	}
	return false
}

func hasValueForField(field model.Field) bool {
	if field.TypeName == "int" || field.TypeName == "string" || field.TypeName == "bool" {
		return true
	}
	return false
}

func valueForField(field model.Field) string {
	if field.TypeName == "int" {
		if field.IsSlice {
			return "[]int{1,2}"
		}
		return "42"
	}

	if field.TypeName == "string" {
		if field.IsSlice {
			return "[]string{" + fmt.Sprintf("\"Example1%s\"", field.Name) + "," +
				fmt.Sprintf("\"Example1%s\"", field.Name) + "}"
		}
		return fmt.Sprintf("\"Example3%s\"", field.Name)
	}

	if field.TypeName == "bool" {
		return "true"
	}
	return ""
}
