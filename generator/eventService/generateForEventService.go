package eventService

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"text/template"
	"unicode"

	"github.com/MarcGrol/golangAnnotations/generator"
	"github.com/MarcGrol/golangAnnotations/generator/annotation"
	"github.com/MarcGrol/golangAnnotations/generator/eventService/eventServiceAnnotation"
	"github.com/MarcGrol/golangAnnotations/generator/generationUtil"
	"github.com/MarcGrol/golangAnnotations/model"
)

type Generator struct {
}

func NewGenerator() generator.Generator {
	return &Generator{}
}

func (eg *Generator) GetAnnotations() []annotation.AnnotationDescriptor {
	return eventServiceAnnotation.Get()
}

func (eg *Generator) Generate(inputDir string, parsedSource model.ParsedSources) error {
	return generate(inputDir, parsedSource.Structs)
}

type templateData struct {
	PackageName string
	Services    []model.Struct
}

func generate(inputDir string, structs []model.Struct) error {

	packageName, err := generationUtil.GetPackageNameForStructs(structs)
	if packageName == "" || err != nil {
		return err
	}
	targetDir, err := generationUtil.DetermineTargetPath(inputDir, packageName)
	if err != nil {
		return err
	}

	eventServices := make([]model.Struct, 0)
	for _, service := range structs {
		if IsEventService(service) {
			eventServices = append(eventServices, service)
		}
	}

	if len(eventServices) == 0 {
		return nil
	}

	data := templateData{
		PackageName: packageName,
		Services:    eventServices,
	}
	return doGenerate(targetDir, packageName, data)
}

func doGenerate(targetDir, packageName string, data templateData) error {
	err := generationUtil.Generate(generationUtil.Info{
		Src:            packageName,
		TargetFilename: generationUtil.Prefixed(fmt.Sprintf("%s/eventHandler.go", targetDir)),
		TemplateName:   "event-handlers",
		TemplateString: handlersTemplate,
		FuncMap:        customTemplateFuncs,
		Data:           data,
	})
	if err != nil {
		log.Fatalf("Error generating handlers for event-services in package %s: %s", packageName, err)
		return err
	}

	for _, eventService := range data.Services {
		if !IsEventServiceNoTest(eventService) {
			err = generationUtil.Generate(generationUtil.Info{
				Src:            packageName,
				TargetFilename: generationUtil.Prefixed(fmt.Sprintf("%s/eventHandlerHelpers_test.go", targetDir)),
				TemplateName:   "test-handlers",
				TemplateString: testHandlersTemplate,
				FuncMap:        customTemplateFuncs,
				Data:           data,
			})
			if err != nil {
				log.Fatalf("Error generating test-handlers for event-services in package %s: %s", packageName, err)
				return err
			}
			break
		}
	}

	return nil
}

var customTemplateFuncs = template.FuncMap{
	"IsEventService":                  IsEventService,
	"IsEventServiceNoTest":            IsEventServiceNoTest,
	"IsEventOperation":                IsEventOperation,
	"GetInputArgType":                 GetInputArgType,
	"GetFullEventNames":               GetFullEventNames,
	"GetInputArgPackage":              GetInputArgPackage,
	"GetEventServiceSelfName":         GetEventServiceSelfName,
	"GetEventServiceTopics":           GetEventServiceTopics,
	"GetEventOperationTopic":          GetEventOperationTopic,
	"IsEventOperationDelayed":         IsEventOperationDelayed,
	"IsAnyEventOperationDelayed":      IsAnyEventOperationDelayed,
	"GetEventOperationQueueGroups":    GetEventOperationQueueGroups,
	"GetEventOperationProducesEvents": GetEventOperationProducesEvents,
	"IsEventNotTransient":             IsEventNotTransient,
	"ToFirstUpper":                    ToFirstUpper,
}

func IsEventService(s model.Struct) bool {
	annotations := annotation.NewRegistry(eventServiceAnnotation.Get())
	_, ok := annotations.ResolveAnnotationByName(s.DocLines, eventServiceAnnotation.TypeEventService)
	return ok
}

func IsEventNotTransient(o model.Operation) bool {
	for _, arg := range o.InputArgs {
		if !IsPrimitiveArg(arg) && !isContextArg(arg) && !isRequestContextArg(arg) {
			// TODO MarcGrol: is there a better way to find out of an event can be stored?
			return !strings.Contains(arg.TypeName, "Discovered")
		}
	}
	return false
}

func IsEventServiceNoTest(s model.Struct) bool {
	annotations := annotation.NewRegistry(eventServiceAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(s.DocLines, eventServiceAnnotation.TypeEventService); ok {
		return ann.Attributes[eventServiceAnnotation.ParamNoTest] == "true"
	}
	return false
}

func GetEventServiceSelfName(s model.Struct) string {
	annotations := annotation.NewRegistry(eventServiceAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(s.DocLines, eventServiceAnnotation.TypeEventService); ok {
		return ann.Attributes[eventServiceAnnotation.ParamSelf]
	}
	return ""
}

func GetEventOperationProducesEventsAsSlice(o model.Operation) []string {
	annotations := annotation.NewRegistry(eventServiceAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(o.DocLines, eventServiceAnnotation.TypeEventOperation); ok {
		if attrs, ok := ann.Attributes[eventServiceAnnotation.ParamProducesEvents]; ok {
			eventsProduced := make([]string, 0)
			for _, e := range strings.Split(attrs, ",") {
				evt := strings.TrimSpace(e)
				if evt != "" {
					eventsProduced = append(eventsProduced, evt)
				}
			}
			return eventsProduced
		}
	}
	return []string{}
}

func GetEventOperationProducesEvents(o model.Operation) string {
	return asStringSlice(GetEventOperationProducesEventsAsSlice(o))
}

func asStringSlice(in []string) string {
	adjusted := make([]string, 0)
	for _, i := range in {
		adjusted = append(adjusted, fmt.Sprintf("\"%s\"", i))
	}
	return fmt.Sprintf("[]string{%s}", strings.Join(adjusted, ","))
}

func GetEventServiceTopics(s model.Struct) []string {
	topics := make([]string, 0)
operations:
	for _, o := range s.Operations {
		if IsEventOperation(*o) {
			topic := GetEventOperationTopic(*o)
			for _, t := range topics {
				if t == topic {
					continue operations
				}
			}
			topics = append(topics, topic)
		}
	}
	return topics
}

func GetFullEventNames(s model.Struct) []string {
	eventMap := map[string]bool{}
	for _, o := range s.Operations {
		if IsEventOperation(*o) {
			eventMap[fmt.Sprintf("%sEvents.%sEventName", GetEventOperationTopic(*o), GetInputArgType(*o))] = true
		}
	}

	eventSlice := make([]string, 0)
	for e := range eventMap {
		eventSlice = append(eventSlice, e)
	}
	sort.Strings(eventSlice)
	return eventSlice
}

func IsEventOperation(o model.Operation) bool {
	annotations := annotation.NewRegistry(eventServiceAnnotation.Get())
	_, ok := annotations.ResolveAnnotationByName(o.DocLines, eventServiceAnnotation.TypeEventOperation)
	return ok
}

func GetEventOperationTopic(o model.Operation) string {
	annotations := annotation.NewRegistry(eventServiceAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(o.DocLines, eventServiceAnnotation.TypeEventOperation); ok {
		return ann.Attributes[eventServiceAnnotation.ParamTopic]
	}
	return ""
}

type queueGroup struct {
	Process string
	Events  []string
}

func GetEventOperationQueueGroups(s model.Struct) []queueGroup {
	queueGroups := make([]queueGroup, 0)
operations:
	for _, o := range s.Operations {
		if IsEventOperation(*o) {
			process := GetEventOperationProcess(*o)
			if process != "" {
				aggregate := GetInputArgPackage(*o)
				eventType := GetInputArgType(*o)
				event := fmt.Sprintf("%s.%s", aggregate, eventType)
				for i, group := range queueGroups {
					if group.Process == process {
						queueGroups[i].Events = append(group.Events, event)
						continue operations
					}
				}
				queueGroups = append(queueGroups, queueGroup{Process: process, Events: []string{event}})
			}
		}
	}
	return queueGroups
}

func GetEventOperationProcess(o model.Operation) string {
	annotations := annotation.NewRegistry(eventServiceAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(o.DocLines, eventServiceAnnotation.TypeEventOperation); ok {
		process := ann.Attributes[eventServiceAnnotation.ParamProcess]
		if process != "" {
			return ToFirstUpper(process)
		}
	}
	return "Default"
}

func IsAnyEventOperationDelayed(s model.Struct) bool {
	for _, oper := range s.Operations {
		if IsEventOperationDelayed(*oper) {
			return true
		}
	}
	return false
}

func IsEventOperationDelayed(o model.Operation) bool {
	annotations := annotation.NewRegistry(eventServiceAnnotation.Get())
	if ann, ok := annotations.ResolveAnnotationByName(o.DocLines, eventServiceAnnotation.TypeEventOperation); ok {
		return ann.Attributes[eventServiceAnnotation.ParamDelayed] == "true"
	}
	return false
}

func GetInputArgType(o model.Operation) string {
	for _, arg := range o.InputArgs {
		if !IsPrimitiveArg(arg) && !isContextArg(arg) && !isRequestContextArg(arg) {
			_, argType := arg.SplitTypeName()
			return argType
		}
	}
	return ""
}

func GetInputArgPackage(o model.Operation) string {
	for _, arg := range o.InputArgs {
		if !IsPrimitiveArg(arg) && !isContextArg(arg) && !isRequestContextArg(arg) {
			argPackage, _ := arg.SplitTypeName()
			return argPackage
		}
	}
	return ""
}
func isContextArg(f model.Field) bool {
	return f.TypeName == "context.Context"
}

func isRequestContextArg(f model.Field) bool {
	return f.TypeName == "request.Context"
}

func IsPrimitiveArg(f model.Field) bool {
	return f.IsPrimitive() || f.IsPrimitiveSlice()
}

func ToFirstUpper(in string) string {
	a := []rune(in)
	a[0] = unicode.ToUpper(a[0])
	return string(a)
}
