package event

import (
	"fmt"
	"log"
	"text/template"
	"unicode"

	"github.com/MarcGrol/golangAnnotations/generator"
	"github.com/MarcGrol/golangAnnotations/generator/annotation"
	"github.com/MarcGrol/golangAnnotations/generator/event/eventAnnotation"
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

func NewGenerator() generator.Generator {
	return &Generator{}
}

func (eg *Generator) GetAnnotations() []annotation.AnnotationDescriptor {
	return eventAnnotation.Get()
}

func (eg *Generator) Generate(inputDir string, parsedSource model.ParsedSources) error {
	return generate(inputDir, parsedSource.Structs)
}

type generateContext struct {
	targetDir   string
	packageName string
	structs     []model.Struct
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

	ctx := generateContext{
		targetDir:   targetDir,
		packageName: packageName,
		structs:     structs,
	}

	err = generateAggregates(ctx)
	if err != nil {
		return err
	}

	err = generateWrappers(ctx)
	if err != nil {
		return err
	}

	err = generateEventStore(ctx)
	if err != nil {
		return err
	}

	err = generateEventPublisher(ctx)
	if err != nil {
		return err
	}

	err = generateWrappersTest(ctx)
	if err != nil {
		return err
	}

	err = generateHandlerInterface(ctx)
	if err != nil {
		return err
	}

	return nil
}

func generateAggregates(ctx generateContext) error {

	aggregates := getAggregates(ctx.structs)

	if len(aggregates) == 0 {
		return nil
	}

	err := generationUtil.Generate(generationUtil.Info{
		Src:            ctx.packageName,
		TargetFilename: generationUtil.Prefixed(fmt.Sprintf("%s/aggregates.go", ctx.targetDir)),
		TemplateName:   "aggregates",
		TemplateString: aggregateTemplate,
		FuncMap:        customTemplateFuncs,
		Data: aggregateMap{
			PackageName:  ctx.packageName,
			AggregateMap: aggregates,
		},
	})
	if err != nil {
		log.Fatalf("Error generating aggregates (%s)", err)
		return err
	}
	return nil
}

func getAggregates(structs []model.Struct) map[string]eventMap {
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
	return aggregates
}

func generateWrappers(ctx generateContext) error {

	if !containsAny(ctx.structs, IsEvent) {
		return nil
	}

	err := generationUtil.Generate(generationUtil.Info{
		Src:            ctx.packageName,
		TargetFilename: generationUtil.Prefixed(fmt.Sprintf("%s/wrappers.go", ctx.targetDir)),
		TemplateName:   "wrappers",
		TemplateString: wrappersTemplate,
		FuncMap:        customTemplateFuncs,
		Data: structures{
			PackageName: ctx.packageName,
			Structs:     ctx.structs,
		},
	})
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

func generateEventStore(ctx generateContext) error {

	if !containsAny(ctx.structs, IsPersistentEvent) {
		return nil
	}

	err := generationUtil.Generate(generationUtil.Info{
		Src:            ctx.packageName,
		TargetFilename: generationUtil.Prefixed(fmt.Sprintf("%s/../%sStore/%sStore.go", ctx.targetDir, ctx.packageName, ctx.packageName)),
		TemplateName:   "event-store",
		TemplateString: eventStoreTemplate,
		FuncMap:        customTemplateFuncs,
		Data: structures{
			PackageName: ctx.packageName,
			Structs:     ctx.structs,
		},
	})
	if err != nil {
		log.Fatalf("Error generating event-store for structures (%s)", err)
		return err
	}
	return nil
}

func generateEventPublisher(ctx generateContext) error {

	if !containsAny(ctx.structs, isTransient) {
		return nil
	}

	err := generationUtil.Generate(generationUtil.Info{
		Src:            ctx.packageName,
		TargetFilename: generationUtil.Prefixed(fmt.Sprintf("%s/../%sPublisher/%sPublisher.go", ctx.targetDir, ctx.packageName, ctx.packageName)),
		TemplateName:   "event-publisher",
		TemplateString: eventPublisherTemplate,
		FuncMap:        customTemplateFuncs,
		Data: structures{
			PackageName: ctx.packageName,
			Structs:     ctx.structs,
		},
	})
	if err != nil {
		log.Fatalf("Error generating event-publisher for structures (%s)", err)
		return err
	}
	return nil
}

func generateWrappersTest(ctx generateContext) error {

	if !containsAny(ctx.structs, IsEvent) {
		return nil
	}

	err := generationUtil.Generate(generationUtil.Info{
		Src:            ctx.packageName,
		TargetFilename: generationUtil.Prefixed(fmt.Sprintf("%s/wrappers_test.go", ctx.targetDir)),
		TemplateName:   "wrappers-test",
		TemplateString: wrappersTestTemplate,
		FuncMap:        customTemplateFuncs,
		Data: structures{
			PackageName: ctx.packageName,
			Structs:     ctx.structs,
		},
	})
	if err != nil {
		log.Fatalf("Error generating wrappers-test for structures (%s)", err)
		return err
	}
	return nil
}

func generateHandlerInterface(ctx generateContext) error {

	if !containsAny(ctx.structs, IsEvent) {
		return nil
	}

	err := generationUtil.Generate(generationUtil.Info{
		Src:            ctx.packageName,
		TargetFilename: generationUtil.Prefixed(fmt.Sprintf("%s/interface.go", ctx.targetDir)),
		TemplateName:   "interface",
		TemplateString: interfaceTemplate,
		FuncMap:        customTemplateFuncs,
		Data: structures{
			PackageName: ctx.packageName,
			Structs:     ctx.structs,
		},
	})
	if err != nil {
		log.Fatalf("Error generating interface for event-handlers (%s)", err)
		return err
	}
	return nil
}

var customTemplateFuncs = template.FuncMap{
	"GetEvents":                 GetEvents,
	"IsEvent":                   IsEvent,
	"IsRootEvent":               IsRootEvent,
	"IsPersistentEvent":         IsPersistentEvent,
	"IsTransientEvent":          IsTransientEvent,
	"GetAggregateName":          GetAggregateName,
	"GetAggregateNameLowerCase": GetAggregateNameLowerCase,
	"HasValueForField":          hasValueForField,
	"ValueForField":             valueForField,
}

func GetEvents(thecontext structures) []model.Struct {
	eventsOnly := make([]model.Struct, 0)
	for _, s := range thecontext.Structs {
		if IsEvent(s) {
			eventsOnly = append(eventsOnly, s)
		}
	}
	return eventsOnly
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

func GetAggregateNameLowerCase(s model.Struct) string {
	return toFirstLower(GetAggregateName(s))
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
	if field.IsPrimitive() || field.IsPrimitiveSlice() {
		return true
	}
	return false
}

func valueForField(field model.Field) string {
	if field.IsInt() || field.IsIntSlice() {
		return valueForIntField(field)
	}

	if field.IsString() || field.IsStringSlice() {
		return valueForStringField(field)
	}

	if field.IsBool() || field.IsBoolSlice() {
		return valueForBoolField(field)
	}
	return ""
}

func valueForIntField(field model.Field) string {
	if field.IsSlice() {
		return "[]int{1,2}"
	}
	return "42"
}

func valueForStringField(field model.Field) string {
	if field.IsSlice() {
		return "[]string{" + fmt.Sprintf("\"Example1%s\"", field.Name) + "," +
			fmt.Sprintf("\"Example1%s\"", field.Name) + "}"
	}
	return fmt.Sprintf("\"Example3%s\"", field.Name)
}

func valueForBoolField(field model.Field) string {
	return "true"
}

func toFirstLower(in string) string {
	a := []rune(in)
	a[0] = unicode.ToLower(a[0])
	return string(a)
}
