package eventService

const handlersTemplate = `// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
		router.HandleFunc("/tasks/{{ $serviceName }}/{{.}}/{eventTypeName}", es.handleHTTPBackgroundEvent()).Methods("POST")
	}
	{{end -}}
}

func (es *{{$eventServiceName}}) enqueueEventToBackground(c context.Context, rc request.Context, topic string, envlp envelope.Envelope) error {
	const subscriber = "{{GetEventServiceSelfName .}}"
	switch envlp.EventTypeName {
		case {{range $idxOper, $evtName := GetFullEventNames .}}{{if $idxOper}}, {{end -}}{{$evtName}}{{end -}}:

			taskURL := fmt.Sprintf("/tasks/{{GetEventServiceSelfName .}}/%s/%s", topic, envlp.EventTypeName)

			asJSON, err := json.Marshal(envlp)
			if err != nil {
				msg := fmt.Sprintf("Error marshalling payload for url '%s'", taskURL)
				myerrorhandling.HandleEventError(c, rc, topic, envlp, msg, err)
				return err
			}

			task := queue.Task{
				Method:  "POST",
				URL:     taskURL,
				Payload: asJSON,
			}

			{{if IsAnyEventOperationDelayed . -}}
			var delay time.Duration
			var eta time.Time
			switch envlp.EventTypeName {
			{{range $oper := .Operations -}}{{if IsEventOperationDelayed $oper -}}
			case {{GetInputArgPackage $oper}}.{{GetInputArgType $oper}}EventName:
				delay, eta = get{{GetInputArgType $oper}}DelayOrETA(c, rc, envlp)
			{{end -}}{{end -}}
			}
			task.Delay = delay
			task.ETA = eta
			{{end -}}

			err = myqueue.AddTask(c, es.getProcessTypeFor(envlp), task)
			if err != nil {
				msg := fmt.Sprintf("Error enqueuing task to url '%s'", taskURL)
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
		case {{range $idx, $event := $queueGroup.Events -}}
		{{if $idx}}, {{end}}{{$event}}EventName{{end}}:
			return myqueue.ProcessType{{$queueGroup.Process}}
		{{end -}}
		default:
			return myqueue.ProcessTypeDefault
	}
}

func (es *{{$eventServiceName}}) handleHTTPBackgroundEvent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := ctx.New.CreateContext(r)
		rc := request.NewMinimalContext(c, r)

		retryCount, err := strconv.Atoi(r.Header.Get("X-AppEngine-TaskRetryCount"))
		if err != nil {
			mylog.New().Error(c, rc, "Error parsing 'X-AppEngine-TaskRetryCount': %s", err)
		}

		if retryCount > 0 && !environ.GetEnvironment(c).RetryFailedEvents(c) {
			mylog.New().Info(c, "Abort retry scheme after %d rertries because of env-setting", retryCount)
			return
		}


		// read and parse request body
		var envlp envelope.Envelope
		err = json.NewDecoder(r.Body).Decode(&envlp)
		if err != nil {
			mylog.New().Error(c, rc, "Error parsing request body (retry-count:%d): %s", retryCount, err)
			errorh.HandleHTTPError(c, rc, errorh.NewInvalidInputErrorf(1, "Error parsing request body (retry-count:%d): %s", retryCount, err), w, r)
			return
		}

		rc.Set(
			request.SessionUID(envlp.SessionUID),
			request.RequestUID(envlp.UUID), // pas a stable identifier that makes writing of resulting events idempotent
			request.TaskRetryCount(retryCount),
		)

		err = es.handleEvent(c, rc, envlp.AggregateName, envlp)
		if err != nil {
			errorh.HandleHTTPError(c, rc, err, w, r)
			return
		}
	}
}

func (es *{{$eventServiceName}}) handleEvent(c context.Context, rc request.Context, topic string, envlp envelope.Envelope) error {
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
					myerrorhandling.HandleEventClearError(c, rc, topic, envlp, fmt.Sprintf("As subscriber '%s': Retry %d of '%s' succeeded", subscriber, rc.GetTaskRetryCount(), envlp.NiceName()))
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
