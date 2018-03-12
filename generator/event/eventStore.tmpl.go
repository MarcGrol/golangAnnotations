package event

const eventStoreTemplate = `// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}Store

import "golang.org/x/net/context"

{{range .Structs -}}

    {{if IsPersistentEvent . -}}

func StoreAndApplyEvent{{.Name}}(c context.Context, rc request.Context, aggregateRoot {{.PackageName}}.{{GetAggregateName .}}Aggregate, event {{.PackageName}}.{{.Name}}) error {
        err := StoreEvent{{.Name}}(c, rc, &event)
        if err == nil {
            aggregateRoot.Apply{{.Name}}(c, event)
        }
        return err
}

// StoreEvent{{.Name}} is used to store event of type {{.Name}}
func StoreEvent{{.Name}}(c context.Context, rc request.Context, event *{{.PackageName}}.{{.Name}}) error {
    envelope, err := event.Wrap(rc)
    if err != nil {
        return errorh.NewInternalErrorf(0, "Error wrapping %s event %s: %s", envelope.EventTypeName, event.GetUID(), err)
    }

    err = store.Put(c, rc, envelope)
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
