package event

const wrappersTestTemplate = `// +build !appengine

// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

import (
    "reflect"
    "testing"
    "time"
    "github.com/stretchr/testify/assert"
)

{{range .Structs -}}
    {{if IsEvent . -}}

func Test{{.Name}}Wrapper(t *testing.T) {
    defer mytime.SetDefaultNow()
	defer myuuid.SetDefaults()

    mytime.SetMockNow()
	myuuid.SetMockV1({{GetAggregateName . }}AggregateName, "1234321")

    event := {{.Name}}{
        {{range .Fields -}}
			{{if HasValueForField . -}}
				{{.Name}}: {{ValueForField .}},
			{{end -}}
		{{end -}}
    }
    wrapped, err := event.Wrap(rest.Credentials{SessionUID:"test_session"})
    assert.NoError(t, err)
    assert.True(t, Is{{.Name}}(wrapped))
    assert.Equal(t, {{GetAggregateName . }}AggregateName, wrapped.AggregateName)
    assert.Equal(t, {{.Name}}EventName, wrapped.EventTypeName)
    //	assert.Equal(t, "UID_{{.Name}}", wrapped.AggregateUID)
    assert.Equal(t, "test_session", wrapped.SessionUID)
    assert.Equal(t, "1234321", wrapped.UUID)
    assert.Equal(t, "2016-02-27T00:00:00+01:00", wrapped.Timestamp.Format(time.RFC3339))
    assert.Equal(t, int64(0), wrapped.SequenceNumber)
    again, ok := GetIfIs{{.Name}}(wrapped)
    assert.True(t, ok)
    assert.NotNil(t,again)
    reflect.DeepEqual(event, *again)
}
        {{end -}}
{{end -}}
`
