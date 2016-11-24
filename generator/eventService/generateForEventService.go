package eventService

import (
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/MarcGrol/golangAnnotations/annotation"
	"github.com/MarcGrol/golangAnnotations/generator/generationUtil"
	"github.com/MarcGrol/golangAnnotations/model"
	"github.com/MarcGrol/golangAnnotations/generator/eventService/eventServiceAnnotation"
)

func Generate(inputDir string, parsedSource model.ParsedSources) error {
	return generate(inputDir, parsedSource.Structs)
}

func generate(inputDir string, structs []model.Struct) error {
	eventServiceAnnotation.Register()

	packageName, err := generationUtil.GetPackageName(structs)
	if err != nil {
		return err
	}
	targetDir, err := generationUtil.DetermineTargetPath(inputDir, packageName)
	if err != nil {
		return err
	}
	for _, service := range structs {
		log.Printf("Found service: %s", service.Name)
		if IsEventService(service) {
			log.Printf("%s is event-service", service.Name)

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
	"IsEventService":           IsEventService,
	"IsEventOperation":         IsEventOperation,
	"GetInputArgType":         GetInputArgType,
	"GetInputArgName":         GetInputArgName,
	"GetInputParamString":     GetInputParamString,
	"GetEventServiceAggregates":GetEventServiceAggregates,
	"GetEventServiceSelfAggregate":GetEventServiceSelfAggregate,
}

func IsEventService(s model.Struct) bool {
	annotation, ok := annotation.ResolveAnnotations(s.DocLines)
	if !ok || annotation.Name != "EventService" {
		return false
	}
	return ok
}


func GetEventServiceSelfAggregate( s model.Struct) string {
	val, ok := annotation.ResolveAnnotations(s.DocLines)
	if ok {
		selfString := val.Attributes["self"]
		return selfString
	}
	return ""
}

func GetEventServiceAggregates(s model.Struct) []string {
	val, ok := annotation.ResolveAnnotations(s.DocLines)
	if ok {
		aggregateString, found := val.Attributes["aggregates"]
		if found {
			return strings.Split(aggregateString, ",")
		}
	}
	return []string{}
}

func IsEventOperation(o model.Operation) bool {
	annotation, ok := annotation.ResolveAnnotations(o.DocLines)
	if !ok || annotation.Name != "EventOperation" {
		return false
	}
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

func GetInputArgName(o model.Operation) string {
	for _, arg := range o.InputArgs {
		if arg.TypeName != "int" && arg.TypeName != "string" && arg.TypeName != "context.Context" {
			return arg.Name
		}
	}
	return ""
}

func GetInputParamString(o model.Operation) string {
	args := []string{}
	for _, arg := range o.InputArgs {
		args = append(args, arg.Name)
	}
	return strings.Join(args, ",")
}

var handlersTemplate string = `
// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

import (
	"golang.org/x/net/context"
	"github.com/Duxxie/platform/backend/lib/bus"
	"github.com/Duxxie/platform/backend/lib/events"
	"github.com/Duxxie/platform/backend/lib/logging"
)

{{ $structName := .Name }}

const (
	subscriber = "{{GetEventServiceSelfAggregate .}}"
)

func init() {

	{{range GetEventServiceAggregates .}}
	{
	    topic := "{{.}}"
	    bus.Subscribe(topic, subscriber, handleEvent)
	}
	{{end}}
}

func handleEvent(c context.Context, topic string, envelope events.Envelope) {
    es := &{{$structName}}{}

	logging.New().Info(c, "As %s: received %s event %s.%s on topic %s",
		subscriber, envelope.EventTypeName, envelope.AggregateName, envelope.AggregateUID, topic)

    {{range $idxOper, $oper := .Operations}}

	{{if IsEventOperation $oper}}
	{
	    event, found := events.GetIfIs{{GetInputArgType $oper}}(&envelope)
	    if found {
		    es.{{$oper.Name}}(c, envelope.SessionUID, *event)
	    }
	}

	{{end}}

{{end}}

}
`
