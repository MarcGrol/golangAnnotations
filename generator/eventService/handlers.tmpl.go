package eventService

const handlersTemplate = `// Generated automatically by golangAnnotations: do not edit manually

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

{{range $idxService, $service := .Services -}}

{{ $eventServiceName := .Name -}}

func (es *{{$eventServiceName}}) SubscribeToEvents(router *mux.Router) {

	const subscriber = "{{GetEventServiceSelfName .}}"
	{{ $serviceName := GetEventServiceSelfName $service }}
	{{range GetEventServiceTopics . -}}
	{
		// Subscribe to topic "{{.}}"
	    bus.Subscribe("{{.}}", subscriber, es.handleEvent)
		{{if IsAsync $service -}}
			router.HandleFunc("/tasks/{{ $serviceName }}/{{.}}/{eventTypeName}", es.httpHandleEventAsync()).Methods("POST")
		{{end -}}
	}
	{{end -}}
}

{{if IsAsync . -}}

func (es *{{$eventServiceName}}) getProcessTypeFor(envlp envelope.Envelope) myqueue.ProcessType {
	switch envlp.EventTypeName {
		{{range $queueGroup := (GetEventOperationQueueGroups .) -}}
		case  {{range $idx, $event := $queueGroup.Events -}}
		{{if $idx}},{{end}}{{$event}}EventName{{end}}:
			return myqueue.ProcessType{{$queueGroup.Process}}
		{{end -}}
		default: return myqueue.ProcessTypeDefault
	}
}

func (es *{{$eventServiceName}}) handleEvent(c context.Context, credentials rest.Credentials, topic string, envlp envelope.Envelope) {
	switch envlp.EventTypeName {
		case {{range $idxOper, $oper := .Operations -}}
			{{if IsEventOperation $oper -}}
				{{if $idxOper}},{{end -}}{{GetInputArgPackage $oper}}.{{GetInputArgType $oper}}EventName{{end -}}
			{{end -}}:

			var delay time.Duration = 0
			{{if IsAnyEventOperationDelayed . -}}
			switch envlp.EventTypeName {
			{{range $oper := .Operations -}}{{if IsEventOperationDelayed $oper -}}
			case {{GetInputArgPackage $oper}}.{{GetInputArgType $oper}}EventName:
				delay = {{GetEventOperationDelay $oper}} * time.Second
			{{end -}}{{end -}}
			}	
			{{end}}

			taskUrl := fmt.Sprintf("/tasks/{{GetEventServiceSelfName .}}/%s/%s", topic, envlp.EventTypeName)

			asJson, err := json.Marshal(envlp)
			if err != nil {
				msg := fmt.Sprintf("Error marshalling payload for url '%s'", taskUrl)
				myerrorhandling.HandleEventError(c, credentials, topic, envlp, msg, err)
				return
			}

			err = myqueue.AddTask(c, es.getProcessTypeFor(envlp), queue.Task{
				Method:  "POST",
				URL:     taskUrl,
				Payload: asJson,
				Delay:   delay,
			})
			if err != nil {
				msg := fmt.Sprintf("Error enqueuing task to url '%s'", taskUrl)
				myerrorhandling.HandleEventError(c, credentials, topic, envlp, msg, err)
				return
			}
			mylog.New().Debug(c, "Subscribed service %s.%s enqueued task on topic '%s' with event %s %s.%s", 
				"{{.PackageName}}", "{{$eventServiceName}}", topic, envlp.EventTypeName, envlp.AggregateName, envlp.AggregateUID )
	}
}

func (es *{{$eventServiceName}}) httpHandleEventAsync() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := ctx.New.CreateContext(r)

		credentials := rest.Credentials{RequestURI: r.RequestURI}

		// read and parse request body
		var envlp envelope.Envelope
		err := json.NewDecoder(r.Body).Decode(&envlp)
		if err != nil {
			rest.HandleHttpError(c, credentials, errorh.NewInvalidInputErrorf(1, "Error parsing request body: %s", err), w, r)
			return
		}
		credentials.SessionUID = envlp.SessionUID
		es.handleEventAsync(c, credentials, envlp.AggregateName, envlp)
	}
}

func (es *{{$eventServiceName}}) handleEventAsync(c context.Context, credentials rest.Credentials, topic string, envlp envelope.Envelope) {
{{else -}}
func (es *{{$eventServiceName}}) handleEvent(c context.Context, credentials rest.Credentials, topic string, envlp envelope.Envelope) {
{{end -}}
	const subscriber = "{{GetEventServiceSelfName .}}"

    {{range $idxOper, $oper := .Operations -}}
		{{if IsEventOperation $oper -}}
		{
			evt, found := {{GetInputArgPackage $oper}}.GetIfIs{{GetInputArgType $oper}}(&envlp)
			if found {
				mylog.New().Debug(c, "-->> As %s: Start handling '%s' for '%s/%s'",
					subscriber, envlp.EventTypeName, envlp.AggregateName, envlp.AggregateUID)
				err := es.{{$oper.Name}}(c, credentials, *evt)
				if err != nil {
					msg := fmt.Sprintf("Subscriber '%s' failed to handle '%s' for '%s/%s'",
						subscriber, envlp.EventTypeName, envlp.AggregateName, envlp.AggregateUID)
					myerrorhandling.HandleEventError(c, credentials, topic, envlp, msg, err)
				} else {
					mylog.New().Debug(c, "<<--As %s: Successfully handled '%s' for '%s/%s'",
						subscriber, envlp.EventTypeName, envlp.AggregateName, envlp.AggregateUID)
				}
			}
		}
		{{end -}}
	{{end -}}
}
{{end -}}
`
