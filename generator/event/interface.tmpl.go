package event

const interfaceTemplate = `// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

import (
	"golang.org/x/net/context"

    "github.com/Duxxie/platform/backend/lib/request"
)

{{$packageName := .PackageName}}

type Handler interface {
{{range GetEvents . -}}
    On{{.Name}}( c context.Context, rc request.Context, event {{.Name}}) error
{{end -}}
}

/*
// These empty implementations can help to easily detect missing methods

{{range $idx, $event := GetEvents . -}}
{{if eq $idx 0 }} 
func forceImplements{{GetAggregateName $event}}EventHandler( specific *{{GetAggregateNameLowerCase $event}}EventService) {{$packageName}}.Handler {
	return specific
}
{{end}}

func (es *{{GetAggregateNameLowerCase $event}}EventService)On{{$event.Name}}( c context.Context, rc request.Context, event {{$packageName}}.{{$event.Name}}) error {
	return es.on{{$event.Name}}(c, rc, event)
}
{{end -}}
*/

`
