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
		err = generationUtil.GenerateFileFromTemplate(templateData, packageName, "handlers", handlersTemplate, customTemplateFuncs, target)
		if err != nil {
			log.Fatalf("Error generating handlers for event-services in package %s: %s", packageName, err)
			return err
		}
	}
	return nil
}

var customTemplateFuncs = template.FuncMap{
	"IsEventService":          IsEventService,
	"IsAsync":                 IsAsync,
	"IsAdmin":                 IsAdmin,
	"IsEventOperation":        IsEventOperation,
	"GetInputArgType":         GetInputArgType,
	"GetInputArgPackage":      GetInputArgPackage,
	"GetEventServiceSelfName": GetEventServiceSelfName,
	"GetEventServiceTopics":   GetEventServiceTopics,
	"GetEventOperationTopic":  GetEventOperationTopic,
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

func IsAdmin(s model.Struct) bool {
	if ann, ok := annotation.ResolveAnnotationByName(s.DocLines, eventServiceAnnotation.TypeEventService); ok {
		adminString, found := ann.Attributes[eventServiceAnnotation.ParamAdmin]
		if found && adminString == "true" {
			return true
		}
	}
	return false
}

func GetEventServiceSelfName(s model.Struct) string {
	if ann, ok := annotation.ResolveAnnotationByName(s.DocLines, eventServiceAnnotation.TypeEventService); ok {
		return ann.Attributes[eventServiceAnnotation.ParamSelf]
	}
	return ""
}

func GetEventOperationProducesEvents(o model.Operation) []string {
	if ann, ok := annotation.ResolveAnnotationByName(o.DocLines, eventServiceAnnotation.TypeEventOperation); ok {
		if atts, ok := ann.Attributes[eventServiceAnnotation.ParamProducesEvents]; ok {
			eventsProduced := strings.Split(atts, ",")
			for i, r := range eventsProduced {
				eventsProduced[i] = strings.Trim(r, " ")
			}
			return eventsProduced
		}
	}
	return []string{}
}

func GetEventServiceTopics(s model.Struct) []string {
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

func GetInputArgType(o model.Operation) string {
	for _, arg := range o.InputArgs {
		if arg.TypeName != "int" && arg.TypeName != "string" && arg.TypeName != "context.Context" && arg.TypeName != "rest.Credentials" {
			tn := strings.Split(arg.TypeName, ".")
			return tn[len(tn)-1]
		}
	}
	return ""
}

func GetInputArgPackage(o model.Operation) string {
	for _, arg := range o.InputArgs {
		if arg.TypeName != "int" && arg.TypeName != "string" && arg.TypeName != "context.Context" && arg.TypeName != "rest.Credentials" {
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
	"encoding/json"
	"fmt"
	"net/http"
	"golang.org/x/net/context"
	"github.com/MarcGrol/golangAnnotations/generator/rest"
	"github.com/MarcGrol/golangAnnotations/generator/rest/errorh"
	"github.com/gorilla/mux"
)

{{range $idxService, $service := .Services}}

{{ $structName := .Name }}

func (es *{{$structName}}) SubscribeToEvents(router *mux.Router) {

	const subscriber = "{{GetEventServiceSelfName .}}"
	{{ $serviceName := GetEventServiceSelfName $service }}
	{{range GetEventServiceTopics .}}
	{
		// Subscribe to topic "{{.}}"
	    bus.Subscribe("{{.}}", subscriber, es.handleEvent)
		{{if IsAsync $service }}
			router.HandleFunc("/tasks/{{ $serviceName }}/{{.}}/{eventTypeName}", es.httpHandleEventAsync()).Methods("POST")
		{{end}}
	}
	{{end}}
}

{{if IsAsync .}}

func (es *{{$structName}}) handleEvent(c context.Context, topic string, envelope envelope.Envelope) {
	switch envelope.EventTypeName {
	case{{range $idxOper, $oper := .Operations}}{{if IsEventOperation $oper}}{{if $idxOper}},{{end}}"{{GetInputArgType $oper}}"{{end}}{{end}}:

		taskUrl := fmt.Sprintf("/tasks/{{GetEventServiceSelfName .}}/%s/%s", topic, envelope.EventTypeName )
		isAdmin := {{if IsAdmin .}}true{{else}}false{{end}}

		asJson, err := json.Marshal(envelope)
		if err != nil {
			msg := fmt.Sprintf("Error marshalling payload for url '%s'", taskUrl)
			event.HandleEventError(c, isAdmin, topic, envelope, msg, err)
			return
		}

		err = queue.New().Add(c, queue.Task{
			Method:  "POST",
			URL:     taskUrl,
			Payload: asJson,
			AdminTask: isAdmin,
		})
		if err != nil {
			msg := fmt.Sprintf("Error enqueuing task to url '%s'", taskUrl)
			event.HandleEventError(c, isAdmin, topic, envelope, msg, err)
			return
		}
		mylog.New().Info(c, "Enqueued task to url %s", taskUrl)
	}
}

func (es *{{$structName}}) httpHandleEventAsync() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := ctx.New.CreateContext(r)

		// read and parse request body
		var envelope envelope.Envelope
		err := json.NewDecoder(r.Body).Decode(&envelope)
		if err != nil {
			rest.HandleHttpError(c, rest.Credentials{}, errorh.NewInvalidInputErrorf(1, "Error parsing request body: %s", err), w, r)
			return
		}
		es.handleEventAsync(c, envelope.AggregateName, envelope)
	}
}

func (es *{{$structName}}) handleEventAsync(c context.Context, topic string, envelope envelope.Envelope) {
{{else}}
func (es *{{$structName}}) handleEvent(c context.Context, topic string, envelope envelope.Envelope) {
{{end}}
	const subscriber = "{{GetEventServiceSelfName .}}"
	credentials := rest.Credentials{SessionUID: envelope.SessionUID}
	isAdmin := {{if IsAdmin .}}true{{else}}false{{end}}

    {{range $idxOper, $oper := .Operations}}
	{{if IsEventOperation $oper}}
	{
	    evt, found := {{GetInputArgPackage $oper}}.GetIfIs{{GetInputArgType $oper}}(&envelope)
	    if found {
			mylog.New().Debug(c, "-->> As %s: Start handling '%s' for aggregate '%s/%s'",
				subscriber, envelope.EventTypeName, envelope.AggregateName, envelope.AggregateUID)
		    err := es.{{$oper.Name}}(c, credentials, *evt)
		    if err != nil {
				mylog.New().Error(c, "<<-- As %s: Error handling '%s' for aggregate '%s/%s': %s",
					subscriber, envelope.EventTypeName, envelope.AggregateName, envelope.AggregateUID, err)
				msg := fmt.Sprintf("Subscriber '%s' failed to handle event", subscriber)
				event.HandleEventError(c, isAdmin, topic, envelope, msg, err)
			} else {
				mylog.New().Debug(c, "<<--As %s: Successfully handled '%s' for aggregate '%s/%s'",
					subscriber, envelope.EventTypeName, envelope.AggregateName, envelope.AggregateUID)
			}
	    }
	}
	{{end}}
{{end}}
}
{{end}}
`
