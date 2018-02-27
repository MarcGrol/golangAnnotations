package eventService

const testHandlersTemplate = `// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

import (
    "golang.org/x/net/context"
    "github.com/gorilla/mux"
	"github.com/Duxxie/platform/backend/lib/request"
)

{{range $idxService, $service := .Services -}}

   {{if not (IsEventServiceNoTest .) -}}

	   {{ $eventService := . -}}
	   {{ $eventServiceName := .Name -}}

       {{range $idxOper, $oper := .Operations -}}
		   {{if IsEventOperation $oper -}}

func {{$oper.Name}}In{{ToFirstUpper $service.Name}}TestHelper(t *testing.T, c context.Context, rc request.Context, es *{{$eventServiceName}}, event {{GetInputArgPackage $oper}}.{{GetInputArgType $oper}} ) []envelope.Envelope{
	{{if IsEventNotTransient $oper -}}
	{
		err := store.StoreEvent(c, rc, &event)
		if err != nil {
			t.Fatalf("Error storing event %s: %s", "{{GetInputArgPackage $oper}}.{{GetInputArgType $oper}}", err)
		}
	}
	{{end -}}

	envlp, err := event.Wrap(rc)
	if err != nil {
		t.Fatalf("Error wrapping event %s: %s", "{{GetInputArgPackage $oper}}.{{GetInputArgType $oper}}", err)
	}

	eventsBefore := getEvents(c, rc)

	es.handleEvent(c, rc, "{{GetEventOperationTopic .}}", *envlp)

	eventsAfter := getEvents(c, rc)
	delta :=  getEventsDelta(eventsBefore, eventsAfter)
	verifyAllowed(t, {{GetEventOperationProducesEvents $oper}},delta)

	return delta
}
            {{end -}}
        {{end -}}
    {{end -}}
{{end -}}

func getEvents(c context.Context, rc request.Context) []envelope.Envelope {
    eventsBefore := []envelope.Envelope{}
    eventStore.Mocked().IterateAll(c, rc, func(e envelope.Envelope) error {
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

func isAllowed(allowedEventNames []string, envlp envelope.Envelope) bool {
    for _, name := range allowedEventNames {
        if name == fmt.Sprintf("%s.%s", envlp.AggregateName, envlp.EventTypeName) {
            return true
        }
    }
    return false
}
`
