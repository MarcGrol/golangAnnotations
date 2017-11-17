package rest

const testHelpersTemplate = `// +build !appengine

// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "net/http/httptest"
    "net/url"
    "os"
    "sort"
    "strings"
    "testing"

    "golang.org/x/net/context"

    "github.com/MarcGrol/golangAnnotations/generator/rest/errorh"
)

{{ $structName := .Name }}

var (
    logFp *os.File
    setCookieHook = func(r *http.Request, headers map[string]string) {}
    eventsForOperations = map[string]map[string]bool{}
)

func openfile( filename string) *os.File {
    fp, err := os.Create(filename)
    if err != nil {
        log.Fatalf("Error opening rest-dump-file %s: %s", filename, err.Error())
    }
    return fp
}

func TestMain(m *testing.M) {

    dirname := "{{.PackageName}}TestLog"
    if _, err := os.Stat(dirname); os.IsNotExist(err) {
        os.Mkdir(dirname, os.ModePerm)
    }
    logFp = openfile(dirname + "/$testResults.go")
    defer func() {
        logFp.Close()
    }()
    fmt.Fprintf(logFp, "package %s\n\n", dirname )
    fmt.Fprintf(logFp, "// Generated automatically based on running of api-tests\n\n" )
    fmt.Fprintf(logFp, "import \"github.com/MarcGrol/golangAnnotations/generator/rest/testcase\"\n")

    fmt.Fprintf(logFp, "var TestResults = testcase.TestSuiteDescriptor {\n" )
    fmt.Fprintf(logFp, "\tPackage: \"{{.PackageName}}\",\n")
    fmt.Fprintf(logFp, "\tTestCases: []testcase.TestCaseDescriptor{\n")

    beforeAll()

    code := m.Run()

    afterAll()

    fmt.Fprintf(logFp, "},\n" )
    fmt.Fprintf(logFp, "}\n" )

    log.Printf("events-for-operations")
    for operationName, events := range eventsForOperations {
        log.Printf("operation: %s -> \"%s\"", operationName, mapToList(events))
    }

    os.Exit(code)
}

var beforeAll = defaultBeforeAll

func defaultBeforeAll() {
    mytime.SetMockNow()
}

var afterAll = defaultAfterAll

func defaultAfterAll() {
    mytime.SetDefaultNow()
}

func testCase(name string, description string) {
    fmt.Fprintf(logFp, "\t\ttestcase.TestCaseDescriptor{\n")
    fmt.Fprintf(logFp, "\t\tName:\"%s\",\n", name)
    fmt.Fprintf(logFp, "\t\tDescription:\"%s\",\n", description)
}

func testCaseDone() {
    fmt.Fprintf(logFp, "},\n")
}

func logOperationEvents(c context.Context, operationName string, allowedEvents []string) func(t *testing.T,c context.Context) {
    eventsBeforeTest := collectBefore(c)
    return func(t *testing.T, c context.Context) {
        collectDelta(t, c, operationName, eventsBeforeTest, allowedEvents)
    }
}

func collectBefore(c context.Context) []envelope.Envelope {
    fmt.Fprintf(logFp, "\tPreConditions: []string{\n")
    eventsBefore := []envelope.Envelope{}
    eventStore.Mocked().IterateAll(c, credentials, func(e envelope.Envelope) error {
        eventsBefore = append(eventsBefore, e)
        fmt.Fprintf(logFp, "\"%s\",\n", fmt.Sprintf("%s.%s", e.AggregateName, e.EventTypeName))
        return nil
    })
    fmt.Fprintf(logFp, "\t},\n")

    return eventsBefore
}

func collectDelta(t *testing.T, c context.Context, operationName string, eventsBefore []envelope.Envelope, allowedEvents []string) []envelope.Envelope {

    after := []envelope.Envelope{}
    eventStore.Mocked().IterateAll(c, credentials, func(e envelope.Envelope) error {
        after = append(after, e)
        return nil
    })

    events, found := eventsForOperations[operationName]
    if !found {
        events = map[string]bool{}
    }

    fmt.Fprintf(logFp, "\tPostConditions: []string{\n")

    createdDuringTest := after[len(eventsBefore):]
    for _, e := range createdDuringTest {
        eventName := fmt.Sprintf("%s.%s", e.AggregateName, e.EventTypeName)
        events[eventName] = true
        if !isEventAllowed(allowedEvents, eventName) {
            t.Fatalf("Event '%s' is NOT allowed as result of operation '%s' (allowed: %+v)", eventName, operationName, allowedEvents)
        }
        fmt.Fprintf(logFp, "\"%s\",\n", eventName)
    }
    fmt.Fprintf(logFp, "\t},\n")

    eventsForOperations[operationName] = events

    return createdDuringTest
}

func mapToList(in map[string]bool) string {
    out := []string{}
    for e, _ := range in {
        out = append(out, e)
    }
    return strings.Join(out, ",")
}

func isEventAllowed(allowedEventNames []string, anEventName string) bool {
    for _, e := range allowedEventNames {
        if anEventName == e {
            return true
        }
    }
    return false
}

{{range .Operations}}

    {{if IsRestOperation . }}

        type {{.Name}}TestRequest struct {
            Url      string
            Headers  map[string]string
            HasInput bool
            HasForm  bool
            {{if HasInput . }}Body {{GetInputArgType . }}{{end}}
            {{if IsRestOperationForm . }}Form url.Values{{end}}
        }

        type {{.Name}}TestResponse struct {
            StatusCode int
            Headers    map[string][]string
            {{if IsRestOperationJSON . }}
                {{if HasOutput . }}
                    Body {{GetOutputArgType . }}
                {{end}}
            {{else}}
                Recorder *httptest.ResponseRecorder
            {{end}}
            ErrorBody  *errorh.Error
        }


        func {{.Name}}TestHelperWithoutHeaders(t *testing.T, c context.Context, url string {{if IsRestOperationForm . }}, form url.Values{{else if HasInput . }}, input {{GetInputArgType . }} {{end}} )  ({{if IsRestOperationJSON . }}int {{if HasOutput . }},{{GetOutputArgType . }}{{end}},*errorh.Error{{else}}*httptest.ResponseRecorder{{end}}, error) {
            return {{.Name}}TestHelperWithHeaders( t, c, url {{if IsRestOperationForm . }}, form{{else if HasInput . }}, input {{end}}, map[string]string{} )
        }

        func {{.Name}}TestHelperWithHeaders(t *testing.T, c context.Context, url string {{if IsRestOperationForm . }}, form url.Values{{else if HasInput . }}, input {{GetInputArgType . }} {{end}}, headers map[string]string)  ({{if IsRestOperationJSON . }}int {{if HasOutput . }},{{GetOutputArgType . }}{{end}},*errorh.Error{{else}}*httptest.ResponseRecorder{{end}}, error) {
            request := {{.Name}}TestRequest{
                Url:     url,
                Headers: headers,
                {{if HasInput . }}Body: input,{{end}}
                {{if IsRestOperationForm .}}Form: form,{{end}}
            }

            response := {{.Name}}TestHelper(t, c, request)

            return {{if IsRestOperationJSON . }}response.StatusCode, {{if HasOutput . }}response.Body,{{end}} response.ErrorBody,{{else}}response.Recorder,{{end}} nil
        }

        func {{.Name}}TestHelper(t *testing.T, c context.Context, request {{.Name}}TestRequest)  {{.Name}}TestResponse {
            fmt.Fprintf(logFp, "\t\tOperation:\"%s\",\n", "{{.Name}}")
            defer func() {
                fmt.Fprintf(logFp, "\t},\n")
            }()

            testcaseCompletion := logOperationEvents(c,  "{{.Name}}", {{GetRestOperationProducesEvents .}})
            defer testcaseCompletion(t, c)

            recorder := httptest.NewRecorder()

            {{if HasUpload . }}
                {{.Name}}SetUpload(request.Body)
                req, err := http.NewRequest("{{GetRestOperationMethod . }}", request.Url, nil)
            {{else if IsRestOperationForm . }}
                req, err := http.NewRequest("{{GetRestOperationMethod . }}", request.Url, strings.NewReader(request.Form.Encode()))
                req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
            {{else if HasInput . }}
                rb, _ := json.Marshal(request.Body)
                // indent for readability
                var requestBody bytes.Buffer
                json.Indent(&requestBody, rb, "", "\t")
                req, err := http.NewRequest("{{GetRestOperationMethod . }}", request.Url, strings.NewReader(requestBody.String()))
            {{else}}
                req, err := http.NewRequest("{{GetRestOperationMethod . }}", request.Url, nil)
            {{end}}
                if err != nil {
                    t.Fatalf("Error creating new request %s", err)
                }
                req.RequestURI = request.Url
            {{if HasUpload . }}
            {{else if HasInput . }}
                req.Header.Set("Content-type", "application/json")
            {{end}}
            {{if HasOutput . }}
                req.Header.Set("Accept", "application/json")
            {{end}}
            for k, v := range request.Headers {
                req.Header.Set(k, v)
            }
            setCookieHook(req, request.Headers)

            headersToBeSorted := []string{}
            for key, values := range req.Header {
                for _, value := range values {
                    headersToBeSorted = append(headersToBeSorted, fmt.Sprintf("%s:%s", key, value))
                }
            }
            sort.Strings(headersToBeSorted)

            fmt.Fprintf(logFp, "\tRequest: testcase.RequestDescriptor{\n")
            fmt.Fprintf(logFp, "\tMethod:\"%s\",\n", "{{GetRestOperationMethod . }}")
            fmt.Fprintf(logFp, "\tUrl:\"%s\",\n", request.Url)
            fmt.Fprintf(logFp, "\tHeaders: []string{\n")
            for _, h := range headersToBeSorted {
                fmt.Fprintf(logFp, "\"%s\",\n", h)
            }
            fmt.Fprintf(logFp, "\t},\n")

            {{if HasUpload . }}
            {{else if HasInput . }}
                fmt.Fprintf(logFp, "\tBody:\n" )
                fmt.Fprintf(logFp, "{{BackTick}}%s{{BackTick}}", requestBody.String() )
                fmt.Fprintf(logFp, ",\n" )
            {{end}}
            fmt.Fprintf(logFp, "},\n")

            // dump readable request
            //payload, err := httputil.DumpRequest(req, true)

            fmt.Fprintf(logFp, "\tResponse:testcase.ResponseDescriptor{\n")
            defer func() {
                fmt.Fprintf(logFp, "\t},\n")
            }()

            webservice := NewRest{{ToFirstUpper $structName}}()
            webservice.HTTPHandler().ServeHTTP(recorder, req)

            {{if IsRestOperationJSON . }}
                // dump readable response
                var responseBody bytes.Buffer
                json.Indent(&responseBody, recorder.Body.Bytes(), "", "\t")
            {{end}}

            fmt.Fprintf(logFp, "\tStatus:%d,\n", recorder.Code)

            headersToBeSorted = []string{}
            for key, values := range recorder.Header() {
                for _, value := range values {
                    headersToBeSorted = append(headersToBeSorted, fmt.Sprintf("%s:%s", key, value))
                }
            }
            sort.Strings(headersToBeSorted)

            fmt.Fprintf(logFp,"\tHeaders:[]string{\n")
            for _, h := range headersToBeSorted {
                fmt.Fprintf(logFp, "\"%s\",\n", h)
            }
            fmt.Fprintf(logFp, "\t},\n")
            fmt.Fprintf(logFp, "\tBody:\n{{BackTick}}%s{{BackTick}},\n", {{if IsRestOperationJSON . }}responseBody.String(){{else}}recorder.Body.Bytes(){{end}})

            {{if IsRestOperationJSON . }}
                {{if HasOutput . }}
                    if recorder.Code != http.StatusOK {
                        // return error response
                        var errorResponse errorh.Error
                        dec := json.NewDecoder(recorder.Body)
                        err = dec.Decode(&errorResponse)
                        if err != nil {
                            t.Fatalf("Error decoding error response %s", err)
                        }

                        return {{.Name}}TestResponse {
                            StatusCode: recorder.Code,
                            Headers:    recorder.Header(),
                            ErrorBody:  &errorResponse,
                        }
                    }

                    // return success response
                    resp := {{GetOutputArgDeclaration . }}
                    dec := json.NewDecoder(recorder.Body)
                    err = dec.Decode({{GetOutputArgName . }})
                    if err != nil {
                        t.Fatalf("Error decoding response %s", err)
                    }

                    return {{.Name}}TestResponse {
                        StatusCode: recorder.Code,
                        Headers:    recorder.Header(),
                        Body:       resp,
                    }
                {{else}}
                    return {{.Name}}TestResponse {
                        StatusCode: recorder.Code,
                        Headers:    recorder.Header(),
                    }
                {{end}}
            {{else}}
                return {{.Name}}TestResponse {
                    StatusCode: recorder.Code,
                    Headers:    recorder.Header(),
                    Recorder:   recorder,
                }
            {{end}}
        }
    {{end}}
{{end}}
`
