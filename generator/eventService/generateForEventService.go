package eventService

import (
	"fmt"
	"log"
	"strings"
	"text/template"
	"unicode"

	"github.com/MarcGrol/golangAnnotations/annotation"
	"github.com/MarcGrol/golangAnnotations/generator/eventService/eventServiceAnnotation"
	"github.com/MarcGrol/golangAnnotations/generator/generationUtil"
	"github.com/MarcGrol/golangAnnotations/model"
)

func Generate(inputDir string, parsedSource model.ParsedSources) error {
	return generate(inputDir, parsedSource.Structs)
}

func generate(inputDir string, structs []model.Struct) error {
	eventServiceAnnotation.Register()

	packageName, err := generationUtil.GetPackageNameForStructs(structs)
	if err != nil {
		return err
	}
	targetDir, err := generationUtil.DetermineTargetPath(inputDir, packageName)
	if err != nil {
		return err
	}

	eventServices := []model.Struct{}
	for _, service := range structs {
		if IsEventService(service) {
			eventServices = append(eventServices, service)
		}
	}

	templateData := struct {
		PackageName string
		Services    []model.Struct
	}{
		PackageName: packageName,
		Services:    eventServices,
	}

	if len(eventServices) > 0 {
		target := fmt.Sprintf("%s/$eventHandler.go", targetDir)
		err = generationUtil.GenerateFileFromTemplateFile(templateData, packageName, "event-handlers", "generator/eventService/handlers.go.tmpl", customTemplateFuncs, target)
		if err != nil {
			log.Fatalf("Error generating handlers for event-services in package %s: %s", packageName, err)
			return err
		}

		for _, eventService := range eventServices {
			if !isEventServiceNoTest(eventService) {
				target = fmt.Sprintf("%s/$eventHandlerHelpers_test.go", targetDir)
				err = generationUtil.GenerateFileFromTemplateFile(templateData, packageName, "test-handlers", "generator/eventService/testHandlers.go.tmpl", customTemplateFuncs, target)
				if err != nil {
					log.Fatalf("Error generating test-handlers for event-services in package %s: %s", packageName, err)
					return err
				}
				break
			}
		}

	}
	return nil
}

var customTemplateFuncs = template.FuncMap{
	"IsEventService":                  IsEventService,
	"IsAsync":                         IsAsync,
	"IsEventServiceNoTest":            isEventServiceNoTest,
	"IsEventOperation":                IsEventOperation,
	"GetInputArgType":                 GetInputArgType,
	"GetInputArgPackage":              getInputArgPackage,
	"GetEventServiceSelfName":         GetEventServiceSelfName,
	"GetEventServiceTopics":           getEventServiceTopics,
	"GetEventOperationTopic":          GetEventOperationTopic,
	"GetEventOperationQueueGroups":    getEventOperationQueueGroups,
	"GetEventOperationProducesEvents": getEventOperationProducesEvents,
	"IsAsyncAsString":                 isAsyncAsString,
	"IsEventNotTransient":             isEventNotTransient,
	"ToFirstUpper":                    toFirstUpper,
}

func IsEventService(s model.Struct) bool {
	_, ok := annotation.ResolveAnnotationByName(s.DocLines, eventServiceAnnotation.TypeEventService)
	return ok
}

func IsAsync(s model.Struct) bool {
	if ann, ok := annotation.ResolveAnnotationByName(s.DocLines, eventServiceAnnotation.TypeEventService); ok {
		syncString, found := ann.Attributes[eventServiceAnnotation.ParamAsync]
		if found && syncString == "true" {
			return true
		}
	}
	return false
}

func isAsyncAsString(s model.Struct) string {
	if IsAsync(s) {
		return "Async"
	}
	return ""
}

func isEventNotTransient(o model.Operation) bool {
	for _, arg := range o.InputArgs {
		if !IsPrimitiveArg(arg) && !isContextArg(arg) && !isCredentialsArg(arg) {
			// TODO MarcGrol: is there a better way to find out of an event can be stored?
			return !strings.Contains(arg.TypeName, "Discovered")
		}
	}
	return false
}

func isEventServiceNoTest(s model.Struct) bool {
	if ann, ok := annotation.ResolveAnnotationByName(s.DocLines, eventServiceAnnotation.TypeEventService); ok {
		return ann.Attributes[eventServiceAnnotation.ParamNoTest] == "true"
	}
	return false
}

func GetEventServiceSelfName(s model.Struct) string {
	if ann, ok := annotation.ResolveAnnotationByName(s.DocLines, eventServiceAnnotation.TypeEventService); ok {
		return ann.Attributes[eventServiceAnnotation.ParamSelf]
	}
	return ""
}

func getEventOperationProducesEventsAsSlice(o model.Operation) []string {
	if ann, ok := annotation.ResolveAnnotationByName(o.DocLines, eventServiceAnnotation.TypeEventOperation); ok {
		if attrs, ok := ann.Attributes[eventServiceAnnotation.ParamProducesEvents]; ok {
			eventsProduced := []string{}
			for _, evt := range strings.Split(attrs, ",") {
				evt := strings.TrimSpace(evt)
				if evt != "" {
					eventsProduced = append(eventsProduced, evt)
				}
			}
			return eventsProduced
		}
	}
	return []string{}
}

func getEventOperationProducesEvents(o model.Operation) string {
	return asStringSlice(getEventOperationProducesEventsAsSlice(o))
}

func asStringSlice(in []string) string {
	adjusted := []string{}
	for _, i := range in {
		adjusted = append(adjusted, fmt.Sprintf("\"%s\"", i))
	}
	return fmt.Sprintf("[]string{%s}", strings.Join(adjusted, ","))
}

func getEventServiceTopics(s model.Struct) []string {
	topics := []string{}
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

func IsEventOperation(o model.Operation) bool {
	_, ok := annotation.ResolveAnnotationByName(o.DocLines, eventServiceAnnotation.TypeEventOperation)
	return ok
}

func GetEventOperationTopic(o model.Operation) string {
	if ann, ok := annotation.ResolveAnnotationByName(o.DocLines, eventServiceAnnotation.TypeEventOperation); ok {
		return ann.Attributes[eventServiceAnnotation.ParamTopic]
	}
	return ""
}

type queueGroup struct {
	Process string
	Events  []string
}

func getEventOperationQueueGroups(s model.Struct) []queueGroup {
	queueGroups := []queueGroup{}
operations:
	for _, o := range s.Operations {
		if IsEventOperation(*o) {
			process := GetEventOperationProcess(*o)
			if process != "" {
				aggregate := getInputArgPackage(*o)
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
	process := ""
	if ann, ok := annotation.ResolveAnnotationByName(o.DocLines, eventServiceAnnotation.TypeEventOperation); ok {
		process = ann.Attributes[eventServiceAnnotation.ParamProcess]
		if process != "" {
			return toFirstUpper(process)
		}
	}
	return "Default"
}

func GetInputArgType(o model.Operation) string {
	for _, arg := range o.InputArgs {
		if !IsPrimitiveArg(arg) && !isContextArg(arg) && !isCredentialsArg(arg) {
			tn := strings.Split(arg.TypeName, ".")
			return tn[len(tn)-1]
		}
	}
	return ""
}

func getInputArgPackage(o model.Operation) string {
	for _, arg := range o.InputArgs {
		if !IsPrimitiveArg(arg) && !isContextArg(arg) && !isCredentialsArg(arg) {
			tn := strings.Split(arg.TypeName, ".")
			return tn[len(tn)-2]
		}
	}
	return ""
}
func isContextArg(f model.Field) bool {
	return f.TypeName == "context.Context"
}

func isCredentialsArg(f model.Field) bool {
	return f.TypeName == "rest.Credentials"
}

func IsPrimitiveArg(f model.Field) bool {
	return isNumberArg(f) || isStringArg(f)
}

func isNumberArg(f model.Field) bool {
	return f.TypeName == "int"
}

func isStringArg(f model.Field) bool {
	return f.TypeName == "string"
}

func toFirstUpper(in string) string {
	a := []rune(in)
	a[0] = unicode.ToUpper(a[0])
	return string(a)
}
