package event

const eventStoreTemplate = `// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}Store

import (
	"context"

	"cloud.google.com/go/datastore"
)

{{range .Structs -}}

	{{if IsPersistentEvent . -}}

func StoreAndApplyEvent{{.Name}}(c context.Context, rc request.Context, tx *datastore.Transaction, aggregateRoot {{.PackageName}}.{{GetAggregateName .}}Aggregate, evt {{.PackageName}}.{{.Name}}) error {
	err := StoreEvent{{.Name}}(c, rc, tx, &evt)
	if err == nil {
		aggregateRoot.Apply{{.Name}}(c, rc, evt)
	}
	return err
}

// StoreEvent{{.Name}} is used to store event of type {{.Name}}
func StoreEvent{{.Name}}(c context.Context, rc request.Context, tx *datastore.Transaction, evt *{{.PackageName}}.{{.Name}}) error {
	envlp, err := evt.Wrap(rc)
	if err != nil {
		return errorh.NewInternalErrorf(0, "Error wrapping %s event %s: %s", envlp.EventTypeName, evt.GetUID(), err)
	}

	err = store.Put(c, rc, tx, envlp)
	if err != nil {
		return errorh.NewInternalErrorf(0, "Error storing %s event %s: %s", envlp.EventTypeName, evt.GetUID(), err)
	}

	evt.Metadata = eventMetaData.Metadata{
		UUID:          envlp.UUID,
		Timestamp:     envlp.Timestamp.In(mytime.DutchLocation),
		EventTypeName: envlp.EventTypeName,
	}

	return nil
}
	{{end -}}
{{end -}}
`
