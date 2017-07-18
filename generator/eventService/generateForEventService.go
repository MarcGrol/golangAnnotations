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
				target := fmt.Sprintf("%s/$eventHandler.go", targetDir)
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
	"IsAsync":                      IsAsync,
	"IsEventOperation":             IsEventOperation,
	"GetInputArgType":              GetInputArgType,
	"GetInputArgPackage":           GetInputArgPackage,
	"GetEventServiceSubscriptions": GetEventServiceSubscriptions,
	"GetEventServiceSelfName":      GetEventServiceSelfName,
}

func IsEventService(s model.Struct) bool {
	_, ok := annotation.ResolveAnnotationByName(s.DocLines, eventServiceAnnotation.TypeEventService)
	return ok
}

func IsAsync(s model.Struct) bool {
	ann, ok := annotation.ResolveAnnotationByName(s.DocLines, eventServiceAnnotation.TypeEventService)
	if ok {
		syncString, found := ann.Attributes[eventServiceAnnotation.ParamAsync]
		if found && syncString == "true" {
			return true
		}
	}
	return false
}

func GetEventServiceSelfName(s model.Struct) string {
	ann, ok := annotation.ResolveAnnotationByName(s.DocLines, eventServiceAnnotation.TypeEventService)
	if ok {
		return ann.Attributes[eventServiceAnnotation.ParamSelf]
	}
	return ""
}

func GetEventServiceSubscriptions(s model.Struct) []string {
	ann, ok := annotation.ResolveAnnotationByName(s.DocLines, eventServiceAnnotation.TypeEventService)
	if ok {
		aggregateString, found := ann.Attributes[eventServiceAnnotation.ParamSubscriptions]
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
	_, ok := annotation.ResolveAnnotationByName(o.DocLines, eventServiceAnnotation.TypeEventOperation)
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

import "golang.org/x/net/context"

{{ $structName := .Name }}

const subscriber = "{{GetEventServiceSelfName .}}"

func (es *{{$structName}}) SubscribeToEvents(router *mux.Router) {
	{{range GetEventServiceSubscriptions .}}
	{
		// Subscribe to topic "{{.}}"
	    bus.Subscribe("{{.}}", subscriber, es.handleEvent)
	}
	{{end}}

	{{if IsAsync .}}
		router.HandleFunc("/tasks/"+subscriber+"/{aggregateName}/{eventTypeName}", es.httpHandleEventAsync()).Methods("POST")
	{{end}}
}

{{if IsAsync .}}

func (es *{{$structName}}) handleEvent(c context.Context, topic string, envelope events.Envelope) {
	switch envelope.EventTypeName {
	case{{range $idxOper, $oper := .Operations}}{{if $idxOper}},{{end}}"{{GetInputArgType $oper}}"{{end}}:

		taskUrl := fmt.Sprintf("/tasks/%s/%s/%s", subscriber, envelope.AggregateName, envelope.EventTypeName )

		asJson, err := json.Marshal(envelope)
		if err != nil {
			mylog.New().Error(c, "Error marshalling payload for task %s for url %s: %s", envelope.EventTypeName, taskUrl, err)
			return
		}

		err = queue.New().Add(c, queue.Task{
			Method:  "POST",
			URL:     taskUrl,
			Payload: asJson,
		})
		if err != nil {
			mylog.New().Error(c, "Error enqueuing task to url %s: %s", taskUrl, err)
			return
		}
		mylog.New().Info(c, "Enqueued task to url %s", taskUrl)
	}
}

func (es *{{$structName}}) httpHandleEventAsync() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := ctx.New.CreateContext(r)

		// read and parse request body
		var envelope events.Envelope
		err := json.NewDecoder(r.Body).Decode(&envelope)
		if err != nil {
			errorhandling.HandleHttpError(c, errorh.NewInvalidInputErrorf(1, "Error parsing request body: %s", err), w)
			return
		}
		es.handleEventAsync(c, envelope.AggregateName, envelope)
	}
}

func (es *{{$structName}}) handleEventAsync(c context.Context, topic string, envelope events.Envelope) {
{{else}}
func (es *{{$structName}}) handleEvent(c context.Context, topic string, envelope events.Envelope) {
{{end}}

    {{range $idxOper, $oper := .Operations}}

	{{if IsEventOperation $oper}}
	{
	    event, found := {{GetInputArgPackage $oper}}.GetIfIs{{GetInputArgType $oper}}(&envelope)
	    if found {
				mylog.New().Debug(c, "-->> As %s: Start handling %s event %s.%s on topic %s",
						subscriber, envelope.EventTypeName, envelope.AggregateName,
						envelope.AggregateUID, topic)
		    err := es.{{$oper.Name}}(c, envelope.SessionUID, *event)
		    if err != nil {
				mylog.New().Error(c, "<<-- As %s: Error handling %s event %s.%s on topic %s: %s",
						subscriber, envelope.EventTypeName, envelope.AggregateName,
						envelope.AggregateUID, topic, err)
			} else {
				mylog.New().Debug(c, "<<--As %s: Successfully handled %s event %s.%s on topic %s",
						subscriber, envelope.EventTypeName, envelope.AggregateName,
						envelope.AggregateUID, topic)
			}
	    }
	}

	{{end}}

{{end}}

}
`
