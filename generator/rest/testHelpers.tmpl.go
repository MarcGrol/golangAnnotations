package rest

const testHelpersTemplate = `// +build !appengine

// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

import (
    "golang.org/x/net/context"
	"github.com/Duxxie/platform/backend/lib/request"
)

var (
    setCookieHook = func(r *http.Request, headers map[string]string) {}
    beforeAll = defaultBeforeAll
    afterAll = defaultAfterAll
    testSuite = libtest.NewHTTPTestSuite("{{.PackageName}}")
)

func TestMain(m *testing.M) {
    beforeAll()

    code := m.Run()

    afterAll()

    // write details of all test-cases in structured readable format
    testSuite.WriteToMarkdownGoVarFile()

    os.Exit(code)
}

{{ $serviceName := .Name -}}

type testClient struct {
	c        context.Context
	t        *testing.T
	testCase *libtest.HTTPTestCase
}

func newTestClient(ctx context.Context, testingT *testing.T, testCase *libtest.HTTPTestCase) *testClient {
	return &testClient{
		c:        ctx,
		t:        testingT,
		testCase: testCase,
	}
}

{{range .Operations -}}

{{if IsRestOperation . -}}


type {{.Name}}TestRequest struct {
    Url      string
    Headers  map[string]string
    {{if HasInput . }}Body {{GetInputArgType . }}{{end}}
    {{if IsRestOperationForm . }}Form url.Values{{end}}
}

type {{.Name}}TestResponse struct {
    StatusCode int
    HeaderMap  http.Header
    GetCookie  func(string) *http.Cookie
    {{if IsRestOperationJSON . }}
        {{if HasOutput . }}
            Body {{GetOutputArgType . }}
        {{end -}}
    {{else}}
        Recorder *httptest.ResponseRecorder
    {{end -}}
    ErrorBody  *errorh.Error
}


func {{.Name}}TestHelperWithoutHeaders(t *testing.T, c context.Context, tc *libtest.HTTPTestCase, url string {{if IsRestOperationForm . }}, form url.Values{{else if HasInput . }}, input {{GetInputArgType . }} {{end}} )  ({{if IsRestOperationJSON . }}int {{if HasOutput . }},{{GetOutputArgType . }}{{end}},*errorh.Error{{else}}*httptest.ResponseRecorder{{end}}, error) {
    return {{.Name}}TestHelperWithHeaders( t, c, tc, url {{if IsRestOperationForm . }}, form{{else if HasInput . }}, input {{end}}, map[string]string{} )
}

func {{.Name}}TestHelperWithHeaders(t *testing.T, c context.Context, tc *libtest.HTTPTestCase, url string {{if IsRestOperationForm . }}, form url.Values{{else if HasInput . }}, input {{GetInputArgType . }} {{end}}, headers map[string]string)  ({{if IsRestOperationJSON . }}int {{if HasOutput . }},{{GetOutputArgType . }}{{end}},*errorh.Error{{else}}*httptest.ResponseRecorder{{end}}, error) {
    request := {{.Name}}TestRequest{
        Url:     url,
        Headers: headers,
        {{if HasInput . }}Body: input,{{end}}
        {{if IsRestOperationForm .}}Form: form,{{end}}
    }

    response := newTestClient(c, t, tc).{{.Name}}(request)

    return {{if IsRestOperationJSON . }}response.StatusCode, {{if HasOutput . }}response.Body,{{end}} response.ErrorBody,{{else}}response.Recorder,{{end}} nil
}

func (tcl *testClient){{.Name}}(request {{.Name}}TestRequest) {{.Name}}TestResponse {

    var err error

    // add operation specific info to test-case
    tcl.testCase.ForOperationName("{{.Name}}").
       WithAllowedPostConditions({{GetRestOperationProducesEvents .}}).
       WithPreConditions(fetchEvents(tcl.c))

    // called when function terminates
    defer func() {
        // verify post-conditions
        tc, err := tcl.testCase.WithPostConditions(fetchEvents(tcl.c))
        if err != nil {
            tcl.t.Fatalf("Invalid post-conditions: %s", err )
        }
        // add recordings of this test-case to the test-suite
        testSuite.Add(tc)
    }()

    // compose http-request
    var httpReq *http.Request
    {
        var requestPayload []byte
        {{if HasUpload . -}}
            {{.Name}}SetUpload(request.Body)
            httpReq, err = http.NewRequest("{{GetRestOperationMethod . }}", request.Url, nil)
        {{else if IsRestOperationForm . -}}
            httpReq, err = http.NewRequest("{{GetRestOperationMethod . }}", request.Url, strings.NewReader(request.Form.Encode()))
            httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
        {{else if HasInput . -}}
            requestPayload, err = json.MarshalIndent(request.Body, "", "\t")
            if err != nil {
                tcl.t.Fatalf("Error marshalling request: %s", err )
            }
            httpReq, err = http.NewRequest("{{GetRestOperationMethod . }}", request.Url, strings.NewReader(string(requestPayload)))
        {{else -}}
            httpReq, err = http.NewRequest("{{GetRestOperationMethod . }}", request.Url, nil)
        {{end -}}
        if err != nil {
            tcl.t.Fatalf("Error creating http-request: %s", err )
        }
        httpReq.RequestURI = request.Url
        {{if HasUpload . -}}
        {{else if HasInput . -}}
            httpReq.Header.Set("Content-type", "application/json")
        {{end -}}
        {{if HasOutput . -}}
            httpReq.Header.Set("Accept", "application/json")
        {{end -}}
        for k, v := range request.Headers {
            httpReq.Header.Set(k, v)
        }
        setCookieHook(httpReq, request.Headers)

        // record request-part of test-case
        tcl.testCase.WithRequest("{{GetRestOperationMethod . }}", request.Url, httpReq.Header, requestPayload)
    }

    // call server
    httpResp := httptest.NewRecorder()
    {
        // invoke business logic on remote service
        webservice := NewRest{{ToFirstUpper $serviceName}}()
        webservice.HTTPHandler().ServeHTTP(httpResp, httpReq)

        // record responsepart of testcase
        tcl.testCase.WithResponse(httpResp.Code, httpResp.Header() , httpResp.Body.Bytes())
    }


    // handle response
    {
        // read cookies
        requestWithCookies := &http.Request{
            Header: http.Header{"Cookie": httpResp.HeaderMap["Set-Cookie"]},
        }

		getCookie := func (name string) *http.Cookie {
			cookie, err := requestWithCookies.Cookie(name)
			if err != nil {
				tcl.t.Logf("Error reading cookie '%s': %s", name, err)
			}
			return cookie
		}

        {{if IsRestOperationJSON . -}}
            {{if HasOutput . -}}
                if httpResp.Code != http.StatusOK {
                    // return type-strong error response
                    var errorResponse errorh.Error
                    dec := json.NewDecoder(httpResp.Body)
                    err = dec.Decode(&errorResponse)
                    if err != nil {
                        tcl.t.Fatalf("Error unmarshalling error-response: %s", err )
                    }

                    return {{.Name}}TestResponse {
                        StatusCode: httpResp.Code,
                        HeaderMap:  httpResp.HeaderMap,
                        GetCookie:	getCookie,
                        ErrorBody:  &errorResponse,
                    }
                }

                // return type-strong success response
                resp := {{GetOutputArgDeclaration . }}
                dec := json.NewDecoder(httpResp.Body)
                err = dec.Decode({{GetOutputArgName . }})
                if err != nil {
                    tcl.t.Fatalf("Error unmarshalling response: %s", err )
                }

                return {{.Name}}TestResponse {
                    StatusCode: httpResp.Code,
                    HeaderMap:  httpResp.HeaderMap,
                    GetCookie:	getCookie,
                    Body:       resp,
                }
            {{else -}}
                return {{.Name}}TestResponse {
                    StatusCode: httpResp.Code,
                    HeaderMap:  httpResp.HeaderMap,
                    GetCookie:	getCookie,
                }
            {{end -}}
        {{else -}}
            return {{.Name}}TestResponse {
                StatusCode: httpResp.Code,
                HeaderMap:  httpResp.HeaderMap,
                GetCookie:	getCookie,
                Recorder:   httpResp,
            }
        {{end -}}
    }
}
    {{end -}}
{{end -}}

func defaultBeforeAll() {
    mytime.SetMockNow()
}

func defaultAfterAll() {
    mytime.SetDefaultNow()
}

func fetchEvents(c context.Context) []string {
    found := []string{}
    eventStore.Mocked().IterateAll(c, request.NewEmptyContext(), func(e envelope.Envelope) error {
        found = append(found, fmt.Sprintf("%s.%s", e.AggregateName, e.EventTypeName))
        return nil
    })
    return found
}
`
