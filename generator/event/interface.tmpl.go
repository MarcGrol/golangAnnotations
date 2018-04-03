package event

const interfaceTemplate = `// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

import (
	"golang.org/x/net/context"

    "github.com/Duxxie/platform/backend/lib/request"
)

{{$packageName := .PackageName}}

type Handler interface {
{{range .Structs -}}
    {{if IsEvent . -}}
        On{{.Name}}( c context.Context, rc request.Context, event {{.Name}}) error
    {{end -}}
{{end -}}
}

/*
// These empty implementations can help to easily detect missing methods

{{range .Structs -}}
    {{if IsEvent .}}
func (es *{{GetAggregateNameLowerCase .}}EventService)On{{.Name}}( _ context.Context, _ request.Context, _ {{$packageName}}.{{.Name}}) error {
	return nil
}
    {{end -}}
{{end -}}
*/

`
