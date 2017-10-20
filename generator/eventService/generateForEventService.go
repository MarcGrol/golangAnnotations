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
		err = generationUtil.GenerateFileFromTemplate(templateData, packageName, "handlers", handlersTemplate, customTemplateFuncs, target)
		if err != nil {
			log.Fatalf("Error generating handlers for event-services in package %s: %s", packageName, err)
			return err
		}

		for _, eventService := range eventServices {
			if !IsEventServiceNoTest(eventService) {
				target = fmt.Sprintf("%s/$eventHandlerHelpers_test.go", targetDir)
				err = generationUtil.GenerateFileFromTemplate(templateData, packageName, "testHandlers", handlersTestTemplate, customTemplateFuncs, target)
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
	"IsEventServiceNoTest":            IsEventServiceNoTest,
	"IsEventOperation":                IsEventOperation,
	"GetInputArgType":                 GetInputArgType,
	"GetInputArgPackage":              GetInputArgPackage,
	"GetEventServiceSelfName":         GetEventServiceSelfName,
	"GetEventServiceTopics":           GetEventServiceTopics,
	"GetEventOperationTopic":          GetEventOperationTopic,
	"GetEventOperationQueueGroups":    GetEventOperationQueueGroups,
	"GetEventOperationProducesEvents": GetEventOperationProducesEvents,
	"IsAsyncAsString":                 IsAsyncAsString,
	"IsEventNotTransient":             IsEventNotTransient,
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

func IsAsyncAsString(s model.Struct) string {
	if IsAsync(s) {
		return "Async"
	}
	return ""
}

func IsEventNotTransient(o model.Operation) bool {
	for _, arg := range o.InputArgs {
		if !IsPrimitiveArg(arg) && !IsContextArg(arg) && !IsCredentialsArg(arg) {
			// TODO MarcGrol: is there a better way to find out of an event can be stored?
			return !strings.Contains(arg.TypeName, "Discovered")
		}
	}
	return false
}

func IsEventServiceNoTest(s model.Struct) bool {
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

func GetEventOperationProducesEventsAsSlice(o model.Operation) []string {
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

func GetEventOperationProducesEvents(o model.Operation) string {
	return asStringSlice(GetEventOperationProducesEventsAsSlice(o))
}

func asStringSlice(in []string) string {
	adjusted := []string{}
	for _, i := range in {
		adjusted = append(adjusted, fmt.Sprintf("\"%s\"", i))
	}
	return fmt.Sprintf("[]string{%s}", strings.Join(adjusted, ","))
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

type queueGroup struct {
	Process string
	Events  []string
}

func GetEventOperationQueueGroups(s model.Struct) []queueGroup {
	queueGroups := []queueGroup{}
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
		if !IsPrimitiveArg(arg) && !IsContextArg(arg) && !IsCredentialsArg(arg) {
			tn := strings.Split(arg.TypeName, ".")
			return tn[len(tn)-1]
		}
	}
	return ""
}

func GetInputArgPackage(o model.Operation) string {
	for _, arg := range o.InputArgs {
		if !IsPrimitiveArg(arg) && !IsContextArg(arg) && !IsCredentialsArg(arg) {
			tn := strings.Split(arg.TypeName, ".")
			return tn[len(tn)-2]
		}
	}
	return ""
}
func IsContextArg(f model.Field) bool {
	return f.TypeName == "context.Context"
}

func IsCredentialsArg(f model.Field) bool {
	return f.TypeName == "rest.Credentials"
}

func IsPrimitiveArg(f model.Field) bool {
	return IsNumberArg(f) || IsStringArg(f)
}

func IsNumberArg(f model.Field) bool {
	return f.TypeName == "int"
}

func IsStringArg(f model.Field) bool {
	return f.TypeName == "string"
}

func toFirstUpper(in string) string {
	a := []rune(in)
	a[0] = unicode.ToUpper(a[0])
	return string(a)
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
		{{if IsAsync $service }}router.HandleFunc("/tasks/{{ $serviceName }}/{{.}}/{eventTypeName}", es.httpHandleEventAsync()).Methods("POST"){{end}}
	}
	{{end}}
}

{{if IsAsync .}}

func (es *{{$structName}}) getProcessTypeFor(env envelope.Envelope) myqueue.ProcessType {
	switch env.EventTypeName {
	{{range $queueGroup := (GetEventOperationQueueGroups .)}}
	case  {{range $idx, $event := $queueGroup.Events}}{{if $idx}},{{end}}{{$event}}EventName{{end}}:
		return myqueue.ProcessType{{$queueGroup.Process}}
	{{end}}
	default: return myqueue.ProcessTypeDefault
	}
}

func (es *{{$structName}}) handleEvent(c context.Context, credentials rest.Credentials, topic string, env envelope.Envelope) {
	switch env.EventTypeName {
	case {{range $idxOper, $oper := .Operations}}{{if IsEventOperation $oper}}{{if $idxOper}},{{end}}{{GetInputArgPackage $oper}}.{{GetInputArgType $oper}}EventName{{end}}{{end}}:

		taskUrl := fmt.Sprintf("/tasks/{{GetEventServiceSelfName .}}/%s/%s", topic, env.EventTypeName)

		asJson, err := json.Marshal(env)
		if err != nil {
			msg := fmt.Sprintf("Error marshalling payload for url '%s'", taskUrl)
			myerrorhandling.HandleEventError(c, credentials, topic, env, msg, err)
			return
		}

		err = myqueue.AddTask(c, es.getProcessTypeFor(env), queue.Task{
			Method:  "POST",
			URL:     taskUrl,
			Payload: asJson,
		})
		if err != nil {
			msg := fmt.Sprintf("Error enqueuing task to url '%s'", taskUrl)
			myerrorhandling.HandleEventError(c, credentials, topic, env, msg, err)
			return
		}
		mylog.New().Info(c, "Enqueued task to url %s", taskUrl)
	}
}

func (es *{{$structName}}) httpHandleEventAsync() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := ctx.New.CreateContext(r)

		credentials := rest.Credentials{RequestURI: r.RequestURI}

		// read and parse request body
		var env envelope.Envelope
		err := json.NewDecoder(r.Body).Decode(&env)
		if err != nil {
			rest.HandleHttpError(c, credentials, errorh.NewInvalidInputErrorf(1, "Error parsing request body: %s", err), w, r)
			return
		}
		credentials.SessionUID = env.SessionUID
		es.handleEventAsync(c, credentials, env.AggregateName, env)
	}
}

func (es *{{$structName}}) handleEventAsync(c context.Context, credentials rest.Credentials, topic string, env envelope.Envelope) {
{{else}}
func (es *{{$structName}}) handleEvent(c context.Context, credentials rest.Credentials, topic string, env envelope.Envelope) {
{{end}}
	const subscriber = "{{GetEventServiceSelfName .}}"

    {{range $idxOper, $oper := .Operations}}
	{{if IsEventOperation $oper}}
	{
	    evt, found := {{GetInputArgPackage $oper}}.GetIfIs{{GetInputArgType $oper}}(&env)
	    if found {
			mylog.New().Debug(c, "-->> As %s: Start handling '%s' for '%s/%s'",
				subscriber, env.EventTypeName, env.AggregateName, env.AggregateUID)
		    err := es.{{$oper.Name}}(c, credentials, *evt)
		    if err != nil {
				msg := fmt.Sprintf("Subscriber '%s' failed to handle '%s' for '%s/%s'",
					subscriber, env.EventTypeName, env.AggregateName, env.AggregateUID)
				myerrorhandling.HandleEventError(c, credentials, topic, env, msg, err)
			} else {
				mylog.New().Debug(c, "<<--As %s: Successfully handled '%s' for '%s/%s'",
					subscriber, env.EventTypeName, env.AggregateName, env.AggregateUID)
			}
	    }
	}
	{{end}}
{{end}}
}
{{end}}
`

var handlersTestTemplate string = `
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

   {{if not (IsEventServiceNoTest .) }}

   {{ $struct := . }}
   {{ $structName := .Name }}

   {{range $idxOper, $oper := .Operations}}
		{{if IsEventOperation $oper}}

		func {{$oper.Name}}In{{$service.Name}}TestHelper(t *testing.T, c context.Context, creds rest.Credentials, es *{{$structName}}, event {{GetInputArgPackage $oper}}.{{GetInputArgType $oper}} ) []envelope.Envelope{
			{{if IsEventNotTransient $oper}}
			{
				err := store.StoreEvent{{GetInputArgType $oper}}(c, creds, &event)
				if err != nil {
					t.Fatalf("Error storing event %s: %s", "{{GetInputArgPackage $oper}}.{{GetInputArgType $oper}}", err)
				}
			}
			{{end}}

			envlp, err := event.Wrap(creds.SessionUID)
			if err != nil {
				t.Fatalf("Error wrapping event %s: %s", "{{GetInputArgPackage $oper}}.{{GetInputArgType $oper}}", err)
			}

			eventsBefore := getEvents(c, creds)

			es.handleEvent{{IsAsyncAsString $struct}}(c, creds, "caregiver", *envlp)

			eventsAfter := getEvents(c, creds)
			delta :=  getEventsDelta(eventsBefore, eventsAfter)
			verifyAllowed(t, {{GetEventOperationProducesEvents $oper}},delta)

			return delta
		}
		{{end}}

    {{end}}

	{{end}}

{{end}}

func getEvents(c context.Context, creds rest.Credentials) []envelope.Envelope {
	eventsBefore := []envelope.Envelope{}
	eventStore.Mocked().IterateAll(c, creds, func(e envelope.Envelope) error {
		eventsBefore = append(eventsBefore, e)
		return nil
	})
	return eventsBefore
}

func getEventsDelta(before, after []envelope.Envelope) []envelope.Envelope {
	return after[len(before):]
}

func verifyAllowed(t *testing.T, allowedNames []string, delta []envelope.Envelope) {
	for _, e := range delta {
		if !isAllowed(allowedNames, e) {
			t.Fatalf("Event %s.%s is not allowed", e.AggregateName, e.EventTypeName)
		}
	}
}

func isAllowed(allowedEventNames []string, env envelope.Envelope) bool {
	for _, name := range allowedEventNames {
		if name == fmt.Sprintf("%s.%s", env.AggregateName, env.EventTypeName) {
			return true
		}
	}
	return false
}
`
