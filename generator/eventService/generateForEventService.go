package eventService

import (
	"fmt"
	"log"
	"strings"
	"text/template"

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
	for _, service := range structs {
		if IsEventService(service) {
			{
				target := fmt.Sprintf("%s/eventHandler.go", targetDir)
				err = generationUtil.GenerateFileFromTemplate(service, fmt.Sprintf("%s.%s", service.PackageName, service.Name), "handlers", handlersTemplate, customTemplateFuncs, target)
				if err != nil {
					log.Fatalf("Error generating handlers for event-service %s: %s", service.Name, err)
					return err
				}
			}
		}
	}
	return nil
}

var customTemplateFuncs = template.FuncMap{
	"IsEventService":               IsEventService,
	"IsEventOperation":             IsEventOperation,
	"GetInputArgType":              GetInputArgType,
	"GetInputArgPackage":           GetInputArgPackage,
	"GetEventServiceSubscriptions": GetEventServiceSubscriptions,
	"GetEventServiceSelfName":      GetEventServiceSelfName,
}

func IsEventService(s model.Struct) bool {
	_, ok := annotation.ResolveAnnotationByName(s.DocLines, string(eventServiceAnnotation.TypeEventService))
	return ok
}

func GetEventServiceSelfName(s model.Struct) string {
	ann, ok := annotation.ResolveAnnotationByName(s.DocLines, string(eventServiceAnnotation.TypeEventService))
	if ok {
		return ann.Attributes[string(eventServiceAnnotation.ParamSelf)]
	}
	return ""
}

func GetEventServiceSubscriptions(s model.Struct) []string {
	ann, ok := annotation.ResolveAnnotationByName(s.DocLines, string(eventServiceAnnotation.TypeEventService))
	if ok {
		aggregateString, found := ann.Attributes[string(eventServiceAnnotation.ParamSubscriptions)]
		if found {
			splitted := strings.Split(aggregateString, ",")
			result := []string{}
			for _, s := range splitted {
				result = append(result, strings.TrimSpace(s))
			}
			return result
		}
	}
	return []string{}
}

func IsEventOperation(o model.Operation) bool {
	_, ok := annotation.ResolveAnnotationByName(o.DocLines, string(eventServiceAnnotation.TypeEventOperation))
	return ok
}

func GetInputArgType(o model.Operation) string {
	for _, arg := range o.InputArgs {
		if arg.TypeName != "int" && arg.TypeName != "string" && arg.TypeName != "context.Context" {
			tn := strings.Split(arg.TypeName, ".")
			return tn[len(tn)-1]
		}
	}
	return ""
}

func GetInputArgPackage(o model.Operation) string {
	for _, arg := range o.InputArgs {
		if arg.TypeName != "int" && arg.TypeName != "string" && arg.TypeName != "context.Context" {
			tn := strings.Split(arg.TypeName, ".")
			return tn[len(tn)-2]
		}
	}
	return ""
}

var handlersTemplate string = `
// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

import (
    "golang.org/x/net/context"
)

{{ $structName := .Name }}

const (
	subscriber = "{{GetEventServiceSelfName .}}"
)

func SubscribeToEvents() {
	{{range GetEventServiceSubscriptions .}}
	{
		// Subscribe to topic "{{.}}"
	    bus.Subscribe("{{.}}", subscriber, handleEvent)
	}
	{{end}}
}

func handleEvent(c context.Context, topic string, envelope events.Envelope) {
    es := &{{$structName}}{}

    {{range $idxOper, $oper := .Operations}}

	{{if IsEventOperation $oper}}
	{
	    event, found := {{GetInputArgPackage $oper}}.GetIfIs{{GetInputArgType $oper}}(&envelope)
	    if found {
				logging.New().Debug(c, "-->> As %s: Start handling %s event %s.%s on topic %s",
						subscriber, envelope.EventTypeName, envelope.AggregateName,
						envelope.AggregateUID, topic)
		    err := es.{{$oper.Name}}(c, envelope.SessionUID, *event)
		    if err != nil {
				logging.New().Error(c, "<<-- As %s: Error handling %s event %s.%s on topic %s: %s",
						subscriber, envelope.EventTypeName, envelope.AggregateName,
						envelope.AggregateUID, topic, err)
			} else {
				logging.New().Debug(c, "<<--As %s: Successfully handled %s event %s.%s on topic %s",
						subscriber, envelope.EventTypeName, envelope.AggregateName,
						envelope.AggregateUID, topic)
			}
	    }
	}

	{{end}}

{{end}}

}
`
