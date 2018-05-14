package eventService

const handlersTemplate = `// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

import (
	"encoding/json"
	"fmt"
	"net/http"
	"golang.org/x/net/context"
	"github.com/gorilla/mux"
)

{{range $idxService, $service := .Services -}}

{{ $eventServiceName := .Name -}}

func (es *{{$eventServiceName}}) SubscribeToEvents(router *mux.Router) {
	const subscriber = "{{GetEventServiceSelfName .}}"
	{{ $serviceName := GetEventServiceSelfName $service }}
	{{range GetEventServiceTopics . -}}
	{
	    bus.Subscribe("{{.}}", subscriber, es.enqueueEventToBackground)
		router.HandleFunc("/tasks/{{ $serviceName }}/{{.}}/{eventTypeName}", es.handleHttpBackgroundEvent()).Methods("POST")
	}
	{{end -}}
}

func (es *{{$eventServiceName}}) enqueueEventToBackground(c context.Context, rc request.Context, topic string, envlp envelope.Envelope) error{
	const subscriber = "{{GetEventServiceSelfName .}}"
	switch envlp.EventTypeName {
		case {{range $idxOper, $evtName := GetFullEventNames .}}{{if $idxOper}},{{end -}}{{$evtName}}{{end -}}:

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
				myerrorhandling.HandleEventError(c, rc, topic, envlp, msg, err)
				return err
			}

			err = myqueue.AddTask(c, es.getProcessTypeFor(envlp), queue.Task{
				Method:  "POST",
				URL:     taskUrl,
				Payload: asJson,
				Delay:   delay,
			})
			if err != nil {
				msg := fmt.Sprintf("Error enqueuing task to url '%s'", taskUrl)
				myerrorhandling.HandleEventError(c, rc, topic, envlp, msg, err)
				return err
			}

			mylog.New().Debug(c, "Subscriber '%s' enqueued task on topic '%s' with event '%s'", subscriber, topic, envlp.NiceName())

			return nil
	}
	return nil
}


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

func (es *{{$eventServiceName}}) handleHttpBackgroundEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := ctx.New.CreateContext(r)
		rc := request.NewMinimalContext(c,r)

		retryCount, err := strconv.Atoi(r.Header.Get("X-AppEngine-TaskRetryCount"))
		if err != nil {
			mylog.New().Error(c, rc, request"Error parsing 'X-AppEngine-TaskRetryCount': %s", err)
		}

		if retryCount > 0 && !environ.GetEnvironment(c).RetryFailedEvents(c) {
			mylog.New().Info(c, "Abort retry scheme after %d rertries because of env-setting", retryCount)
			return
		}


		// read and parse request body
		var envlp envelope.Envelope
		err = json.NewDecoder(r.Body).Decode(&envlp)
		if err != nil {
			errorh.HandleHttpError(c, rc, errorh.NewInvalidInputErrorf(1, "Error parsing request body (retry-count:%d): %s", retryCount, err), w, r)
			return
		}

		rc.Set(
		    request.SessionUID(envlp.SessionUID),
		    request.RequestUID(envlp.UUID), // pas a stable identifier that makes writing of resulting events idempotent
			request.TaskRetryCount(retryCount),
		) 

		err = es.handleEvent(c, rc, envlp.AggregateName, envlp)
		if err != nil {
			errorh.HandleHttpError(c, rc, err, w, r)
			return
		}
	}
}

func (es *{{$eventServiceName}}) handleEvent(c context.Context, rc request.Context, topic string, envlp envelope.Envelope) error{
	const subscriber = "{{GetEventServiceSelfName .}}"

    {{range $idxOper, $oper := .Operations -}}
		{{if IsEventOperation $oper -}}
		{
			evt, found := {{GetInputArgPackage $oper}}.GetIfIs{{GetInputArgType $oper}}(&envlp)
			if found {
				err := es.{{$oper.Name}}(c, rc, *evt)
				if err != nil {
					msg := fmt.Sprintf("As subscriber '%s': Failed to handle '%s' (retry: %d)", subscriber, envlp.NiceName(), rc.GetTaskRetryCount())
					myerrorhandling.HandleEventError(c, rc, topic, envlp, msg, err)
					return err
				}

				if rc.GetTaskRetryCount() > 0 {
					myerrorhandling.HandleEventClearError(c, rc, topic, envlp, fmt.Sprintf("As subscriber '%s': Retry %d of '%s' succeeded", subscriber, rc.GetTaskRetryCount(), envlp.NiceName() ))
				}

				return nil
			}
		}
		{{end -}}
	{{end -}}
	return nil
}
{{end -}}
`
