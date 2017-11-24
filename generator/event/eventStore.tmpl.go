package event

const eventStoreTemplate = `// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}Store

import (
    "golang.org/x/net/context"
    "github.com/MarcGrol/golangAnnotations/generator/rest"
    "github.com/MarcGrol/golangAnnotations/generator/rest/errorh"
)

var eventStoreInstance eventStore.EventStore

func init() {
	eventStoreInstance = eventStore.New(myalerts.MyAlertHandler)
}

{{range .Structs -}}

    {{if and (IsEvent .) (IsPersistent .) -}}

func StoreAndApplyEvent{{.Name}}(c context.Context, credentials rest.Credentials, aggregateRoot {{.PackageName}}.{{GetAggregateName .}}Aggregate, event {{.PackageName}}.{{.Name}}) error {
        err := StoreEvent{{.Name}}(c, credentials, &event)
        if err == nil {
            aggregateRoot.Apply{{.Name}}(c, event)
        }
        return err
}

// StoreEvent{{.Name}} is used to store event of type {{.Name}}
func StoreEvent{{.Name}}(c context.Context, credentials rest.Credentials, event *{{.PackageName}}.{{.Name}}) error {
    envelope, err := event.Wrap(credentials.SessionUID)
    if err != nil {
        return errorh.NewInternalErrorf(0, "Error wrapping %s event %s: %s", envelope.EventTypeName, event.GetUID(), err)
    }

    err = eventStoreInstance.Put(c, credentials, envelope)
    if err != nil {
        return errorh.NewInternalErrorf(0, "Error storing %s event %s: %s", envelope.EventTypeName, event.GetUID(), err)
    }

    event.Metadata = {{.PackageName}}.Metadata{
        UUID:          envelope.UUID,
        Timestamp:     envelope.Timestamp.In(mytime.DutchLocation),
        EventTypeName: envelope.EventTypeName,
    }

    return nil
}
	{{end -}}
{{end -}}
`
