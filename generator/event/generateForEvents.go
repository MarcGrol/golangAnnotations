package event

import (
	"fmt"
	"log"
	"text/template"

	"github.com/MarcGrol/golangAnnotations/annotation"
	"github.com/MarcGrol/golangAnnotations/generator/event/eventAnnotation"
	"github.com/MarcGrol/golangAnnotations/generator/filegen"
	"github.com/MarcGrol/golangAnnotations/generator/generationUtil"
	"github.com/MarcGrol/golangAnnotations/model"
)

type eventMap struct {
	Events          map[string]event
	IsAnyPersistent bool
}

type event struct {
	Name         string
	IsPersistent bool
}

type aggregateMap struct {
	PackageName  string
	AggregateMap map[string]eventMap
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

func (eg *Generator) Generate(inputDir string, parsedSource model.ParsedSources) error {
	return generate(inputDir, parsedSource.Structs)
}

func generate(inputDir string, structs []model.Struct) error {
	packageName, err := generationUtil.GetPackageNameForStructs(structs)
	if err != nil {
		return err
	}

	targetDir, err := generationUtil.DetermineTargetPath(inputDir, packageName)
	if err != nil {
		return err
	}

	err = generateAggregates(targetDir, packageName, structs)
	if err != nil {
		return err
	}

	err = generateWrappers(targetDir, packageName, structs)
	if err != nil {
		return err
	}

	err = generateEventStore(targetDir, packageName, structs)
	if err != nil {
		return err
	}

	err = generateEventPublisher(targetDir, packageName, structs)
	if err != nil {
		return err
	}

	err = generateWrappersTest(targetDir, packageName, structs)
	if err != nil {
		return err
	}

	return nil
}

func generateAggregates(targetDir, packageName string, structs []model.Struct) error {

	aggregates := make(map[string]eventMap)
	eventCount := 0
	for _, s := range structs {
		if IsEvent(s) {
			events, ok := aggregates[GetAggregateName(s)]
			if !ok {
				events = eventMap{
					Events:          make(map[string]event),
					IsAnyPersistent: false,
				}
			}
			evt := event{
				Name:         s.Name,
				IsPersistent: IsPersistentEvent(s),
			}
			if evt.IsPersistent {
				events.IsAnyPersistent = true
			}
			events.Events[s.Name] = evt
			aggregates[GetAggregateName(s)] = events
			eventCount++
		}
	}

	if eventCount == 0 {
		return nil
	}

	data := aggregateMap{
		PackageName:  packageName,
		AggregateMap: aggregates,
	}

	target := filegen.Prefixed(fmt.Sprintf("%s/aggregates.go", targetDir))
	err := generationUtil.GenerateFileFromTemplate(data, packageName, "aggregates", aggregateTemplate, customTemplateFuncs, target)
	if err != nil {
		log.Fatalf("Error generating aggregates (%s)", err)
		return err
	}
	return nil
}

func generateWrappers(targetDir, packageName string, structs []model.Struct) error {

	if !containsAny(structs, IsEvent) {
		return nil
	}

	data := structures{
		PackageName: packageName,
		Structs:     structs,
	}
	target := filegen.Prefixed(fmt.Sprintf("%s/wrappers.go", targetDir))
	err := generationUtil.GenerateFileFromTemplate(data, packageName, "wrappers", wrappersTemplate, customTemplateFuncs, target)
	if err != nil {
		log.Fatalf("Error generating wrappers for structures (%s)", err)
		return err
	}
	return nil
}

func containsAny(structs []model.Struct, predicate func(_ model.Struct) bool) bool {
	for _, s := range structs {
		if predicate(s) {
			return true
		}
	}
	return false
}

func generateEventStore(targetDir, packageName string, structs []model.Struct) error {

	if !containsAny(structs, IsPersistentEvent) {
		return nil
	}

	data := structures{
		PackageName: packageName,
		Structs:     structs,
	}
	target := filegen.Prefixed(fmt.Sprintf("%s/../store/%sStore/%sStore.go", targetDir, packageName, packageName))
	err := generationUtil.GenerateFileFromTemplate(data, packageName, "event-store", eventStoreTemplate, customTemplateFuncs, target)
	if err != nil {
		log.Fatalf("Error generating event-store for structures (%s)", err)
		return err
	}
	return nil
}

func generateEventPublisher(targetDir, packageName string, structs []model.Struct) error {

	if !containsAny(structs, isTransient) {
		return nil
	}

	data := structures{
		PackageName: packageName,
		Structs:     structs,
	}
	target := filegen.Prefixed(fmt.Sprintf("%s/../publisher/%sPublisher/%sPublisher.go", targetDir, packageName, packageName))
	err := generationUtil.GenerateFileFromTemplate(data, packageName, "event-publisher", eventPublisherTemplate, customTemplateFuncs, target)
	if err != nil {
		log.Fatalf("Error generating event-publisher for structures (%s)", err)
		return err
	}
	return nil
}

func generateWrappersTest(targetDir, packageName string, structs []model.Struct) error {

	if !containsAny(structs, IsEvent) {
		return nil
	}

	data := structures{
		PackageName: packageName,
		Structs:     structs,
	}
	target := filegen.Prefixed(fmt.Sprintf("%s/wrappers_test.go", targetDir))
	err := generationUtil.GenerateFileFromTemplate(data, packageName, "wrappers-test", wrappersTestTemplate, customTemplateFuncs, target)
	if err != nil {
		log.Fatalf("Error generating wrappers-test for structures (%s)", err)
		return err
	}
	return nil
}

var customTemplateFuncs = template.FuncMap{
	"IsEvent":           IsEvent,
	"IsRootEvent":       IsRootEvent,
	"IsPersistentEvent": IsPersistentEvent,
	"IsTransientEvent":  IsTransientEvent,
	"GetAggregateName":  GetAggregateName,
	"HasValueForField":  hasValueForField,
	"ValueForField":     valueForField,
}

func IsEvent(s model.Struct) bool {
	annotations := annotation.NewRegistry(eventAnnotation.Get())
	_, ok := annotations.ResolveAnnotationByName(s.DocLines, eventAnnotation.TypeEvent)
	return ok
}

func GetAggregateName(s model.Struct) string {
	annotations := annotation.NewRegistry(eventAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(s.DocLines, eventAnnotation.TypeEvent); ok {
		return ann.Attributes[eventAnnotation.ParamAggregate]
	}
	return ""
}

func IsRootEvent(s model.Struct) bool {
	annotations := annotation.NewRegistry(eventAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(s.DocLines, eventAnnotation.TypeEvent); ok {
		return ann.Attributes[eventAnnotation.ParamIsRootEvent] == "true"
	}
	return false
}

func IsPersistentEvent(s model.Struct) bool {
	return IsEvent(s) && !isTransient(s)
}

func IsTransientEvent(s model.Struct) bool {
	return IsEvent(s) && isTransient(s)
}

func isTransient(s model.Struct) bool {
	annotations := annotation.NewRegistry(eventAnnotation.Get())
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
