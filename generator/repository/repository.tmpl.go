package repository

const repositoryTemplate = `// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

import (
    "golang.org/x/net/context"

    "github.com/Duxxie/platform/backend/lib/envelope"
    "github.com/MarcGrol/golangAnnotations/generator/rest"
    "github.com/MarcGrol/golangAnnotations/generator/rest/errorh"
)

{{if HasMethodFind . -}}
var Find{{UpperModelName .}}OnUID = DefaultFind{{UpperModelName .}}OnUID

func DefaultFind{{UpperModelName .}}OnUID(c context.Context, credentials rest.Credentials, {{LowerModelName .}}UID string) (*{{ModelPackageName .}}.{{UpperModelName .}}, error) {
    {{LowerModelName .}}, _, err := DoFind{{UpperModelName .}}OnUID(c, credentials, {{LowerModelName .}}UID, envelope.AcceptAll)
    return {{LowerModelName .}}, err
}

{{if HasMethodFilterByEvent . -}}
func Find{{UpperModelName .}}OnUIDAndEvent(c context.Context, credentials rest.Credentials, {{LowerModelName .}}UID string, metadata {{GetPackageName .}}.Metadata) (*{{ModelPackageName .}}.{{UpperModelName .}}, error) {
    {{LowerModelName .}}, _, err := DoFind{{UpperModelName .}}OnUID(c, credentials, {{LowerModelName .}}UID, envelope.FilterByEventUID{EventUID: metadata.UUID})
    return {{LowerModelName .}}, err
}
{{end -}}

{{if HasMethodFilterByMoment . -}}
func Find{{UpperModelName .}}OnUIDAndMoment(c context.Context, credentials rest.Credentials, {{LowerModelName .}}UID string, moment time.Time) (*{{ModelPackageName .}}.{{UpperModelName .}}, error) {
    {{LowerModelName .}}, _, err := DoFind{{UpperModelName .}}OnUID(c, credentials, {{LowerModelName .}}UID, envelope.FilterByMoment{Moment: moment})
    return {{LowerModelName .}}, err
}
{{end -}}

func DoFind{{UpperModelName .}}OnUID(c context.Context, credentials rest.Credentials, {{LowerModelName .}}UID string, envelopeFilter envelope.EnvelopeFilter) (*{{ModelPackageName .}}.{{UpperModelName .}}, []envelope.Envelope, error) {
    envelopes, err := doFind{{UpperAggregateName .}}EnvelopesOnUID(c, credentials, {{LowerModelName .}}UID, envelopeFilter)
    if err != nil {
        return nil, nil, err
    }

    {{LowerModelName .}} := {{ModelPackageName .}}.New{{UpperModelName .}}()
    err = {{GetPackageName .}}.Apply{{UpperAggregateName .}}Events(c, envelopes, {{LowerModelName .}})
    if err != nil {
        return nil, nil, errorh.NewInternalErrorf(0, "Failed to apply %d events for {{LowerModelName .}} with uid %s: %s", len(envelopes), {{LowerModelName .}}UID, err)
    }
    return {{LowerModelName .}}, envelopes, nil
}

func doFind{{UpperAggregateName .}}EnvelopesOnUID(c context.Context, credentials rest.Credentials, {{LowerModelName .}}UID string, envelopeFilter envelope.EnvelopeFilter) ([]envelope.Envelope, error) {
    envelopes, err := eventStoreInstance.Search(c, credentials, {{GetPackageName .}}.{{AggregateNameConst .}}, {{LowerModelName .}}UID)
    if err != nil {
        return nil, errorh.NewInternalErrorf(0, "Failed to fetch events for {{LowerModelName .}} with uid %s: %s", {{LowerModelName .}}UID, err)
    }

    if len(envelopes) == 0 {
        return nil, errorh.NewNotFoundErrorf(0, "{{UpperModelName .}} with uid %s not found", {{LowerModelName .}}UID)
    }

    envelopes = envelopeFilter.FilteredEnvelopes(envelopes)

    return envelopes, nil
}
{{end -}}

{{if HasMethodFindStates . -}}
    func Find{{UpperModelName .}}StatesOnUID(c context.Context, credentials rest.Credentials, {{LowerModelName .}}UID string) ([]{{ModelPackageName .}}.{{UpperModelName .}}, error) {
    envelopes, err := doFind{{UpperModelName .}}EnvelopesOnUID(c, credentials, {{LowerModelName .}}UID, envelope.AcceptAll)
    if err != nil {
        return nil, err
    }

    states := make([]{{ModelPackageName .}}.{{UpperModelName .}}, 0, len(envelopes))
    {{LowerModelName .}} := {{ModelPackageName .}}.New{{UpperModelName .}}()
    for _, envlp := range envelopes {
        err = {{GetPackageName .}}.Apply{{UpperAggregateName .}}Event(c, envlp, {{LowerModelName .}})
        if err != nil {
            return nil, errorh.NewInternalErrorf(0, "Failed to apply '%s' for {{LowerModelName .}} with uid %s: %s", envlp.EventTypeName, {{LowerModelName .}}UID, err)
        }
        states = append(states, *{{LowerModelName .}})
    }
    return states, nil
    }
{{end -}}

{{if HasMethodExists . -}}
func Exists{{UpperModelName .}}OnUID(c context.Context, credentials rest.Credentials, {{LowerModelName .}}UID string) (bool, error) {
    exists, err := eventStoreInstance.Exists(c, credentials, {{GetPackageName .}}.{{AggregateNameConst .}}, {{LowerModelName .}}UID)
    if err != nil {
        return false, errorh.NewInternalErrorf(0, "Failed to fetch events for {{LowerModelName .}} with uid %s: %s", {{LowerModelName .}}UID, err)
    }
    return exists, nil
}
{{end -}}

{{if HasMethodAllAggregateUIDs . -}}
func GetAll{{UpperModelName .}}UIDs(c context.Context, credentials rest.Credentials) ([]string, error) {
    {{LowerModelName .}}UIDs, err := eventStoreInstance.GetAllAggregateUIDs(c, credentials, {{GetPackageName .}}.{{AggregateNameConst .}})
    if err != nil {
        return nil, errorh.NewInternalErrorf(0, "Failed to fetch all {{LowerModelName .}} uids: %s", err)
    }
        return {{LowerModelName .}}UIDs, nil
    }
{{end -}}

{{if HasMethodGetAllAggregates . -}}
func GetAllRecent{{UpperModelName .}}s(c context.Context, credentials rest.Credentials, optOffset time.Time) ([]{{ModelPackageName .}}.{{UpperModelName .}}, error) {
            {{LowerModelName .}}, _, err := DoGetAllRecent{{UpperModelName .}}s(c, credentials, optOffset)
    return {{LowerModelName .}}, err
}

func DoGetAllRecent{{UpperModelName .}}s(c context.Context, credentials rest.Credentials, optOffset time.Time) ([]{{ModelPackageName .}}.{{UpperModelName .}}, map[string][]envelope.Envelope, error) {
            {{LowerModelName .}}Map := map[string][]envelope.Envelope{}
    err := eventStoreInstance.IterateWithOffset(c, credentials, {{GetPackageName .}}.{{AggregateNameConst .}}, optOffset, func(envlp envelope.Envelope) error {
        if envlp.IsRootEvent {
            {{LowerModelName .}}Map[envlp.AggregateUID] = []envelope.Envelope{envlp}
        } else {
            if envelopes, ok := {{LowerModelName .}}Map[envlp.AggregateUID]; ok {
                {{LowerModelName .}}Map[envlp.AggregateUID] = append(envelopes, envlp)
            }
        }
        return nil
    })
    if err != nil {
        return nil, nil, err
    }

    {{LowerModelName .}}s := make([]{{ModelPackageName .}}.{{UpperModelName .}}, 0, len({{LowerModelName .}}Map))
    for _, {{LowerAggregateName .}}Envelopes := range {{LowerModelName .}}Map {
        {{LowerModelName .}} := {{ModelPackageName .}}.New{{UpperModelName .}}()
        {{GetPackageName .}}.Apply{{UpperAggregateName .}}Events(c, {{LowerAggregateName .}}Envelopes, {{LowerModelName .}})
        {{LowerModelName .}}s = append({{LowerModelName .}}s, *{{LowerModelName .}})
    }
    return {{LowerModelName .}}s, {{LowerModelName .}}Map, nil
}
{{end -}}

{{if HasMethodPurgeOnEventUIDs . -}}
    func Purge{{UpperAggregateName .}}EnvelopesOnUID(c context.Context, credentials rest.Credentials, {{LowerModelName .}}UID string, eventUUIDs []string) error {
    return eventStoreInstance.Purge(c, credentials, {{GetPackageName .}}.{{AggregateNameConst .}}, {{LowerModelName .}}UID, eventUUIDs)
}
{{end -}}

{{if HasMethodPurgeOnEventType . -}}
func PurgeAll{{UpperAggregateName .}}EnvelopesOnEventType(c context.Context, credentials rest.Credentials, eventTypeName string) (bool, error) {
    if eventTypeName == "" {
    return false, errorh.NewInvalidInputErrorf(0, "Missing eventTypeName")
    }
    done, err := eventStoreInstance.PurgeAll(c, credentials, {{GetPackageName .}}.{{AggregateNameConst .}}, eventTypeName)
    if err != nil {
    return false, errorh.NewInternalErrorf(0, "Failed to purge all '%s/%s' events: %s", {{GetPackageName .}}.{{AggregateNameConst .}}, eventTypeName, err)
    }
    return done, nil
}
{{end -}}

{{if HasMethodPurgeAll . -}}
func PurgeAll{{UpperAggregateName .}}Envelopes(c context.Context, credentials rest.Credentials) (bool, error) {
    done, err := eventStoreInstance.PurgeAll(c, credentials, {{GetPackageName .}}.{{AggregateNameConst .}}, "")
    if err != nil {
    return false, errorh.NewInternalErrorf(0, "Failed to purge all '%s' events: %s", {{GetPackageName .}}.{{AggregateNameConst .}}, err)
    }
    return done, nil
}
{{end -}}
`
