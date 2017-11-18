package rest

const testHelpersTemplate = `// +build !appengine

// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

import (
    "golang.org/x/net/context"
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

{{ $serviceName := .Name }}

{{range .Operations}}

{{if IsRestOperation . }}
func {{.Name}}TestHelper(t *testing.T, c context.Context, tc *libtest.HTTPTestCase, url string {{if IsRestOperationForm . }}, form url.Values{{else if HasInput . }}, input {{GetInputArgType . }} {{end}} )  ({{if IsRestOperationJSON . }}int {{if HasOutput . }},{{GetOutputArgType . }}{{end}},*errorh.Error{{else}}*httptest.ResponseRecorder{{end}},error) {
    return {{.Name}}TestHelperWithHeaders( t, c, tc, url {{if IsRestOperationForm . }}, form{{else if HasInput . }}, input {{end}}, map[string]string{} )
}

func {{.Name}}TestHelperWithHeaders(t *testing.T, c context.Context,  tc *libtest.HTTPTestCase, url string {{if IsRestOperationForm . }}, form url.Values{{else if HasInput . }}, input {{GetInputArgType . }} {{end}}, headers map[string]string)  ({{if IsRestOperationJSON . }}int {{if HasOutput . }},{{GetOutputArgType . }}{{end}},*errorh.Error{{else}}*httptest.ResponseRecorder{{end}},error) {
	var err error

	// add operation specific info to test-case
	tc.ForOperationName("{{.Name}}").
       WithAllowedPostConditions({{GetRestOperationProducesEvents .}}).
	   WithPreConditions(fetchEvents(c))

	// called when function terminates
	defer func() {
        // verify post-conditions
		tc, err := tc.WithPostConditions(fetchEvents(c))
		if err != nil {
			t.Fatalf("Invalid post-condions: %s", err )
		}
		// add recordings of this test-case to the test-suite
        testSuite.Add(tc)
	}()

    // compose http-request
	var httpReq *http.Request
	{
		{{if HasUpload . }}
			{{.Name}}SetUpload(input)
			httpReq, err = http.NewRequest("{{GetRestOperationMethod . }}", url, nil)
		{{else if IsRestOperationForm . }}
			httpReq, err = http.NewRequest("{{GetRestOperationMethod . }}", url, strings.NewReader(form.Encode()))
			httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		{{else if HasInput . }}
			requestJson, err := json.MarshalIndent(input, "", "\t")
			if err != nil {
				t.Fatalf("Error marshalling request: %s", err )
			}
			httpReq, err = http.NewRequest("{{GetRestOperationMethod . }}", url, strings.NewReader(string(requestJson)))
		{{else}}
			httpReq, err = http.NewRequest("{{GetRestOperationMethod . }}", url, nil)
		{{end}}
			if err != nil {
				{{if IsRestOperationJSON . }}
					{{if HasOutput . }} return 0, nil, nil, err{{else}}return 0, nil, err{{end}}
				{{else}}return nil, err{{end}}
			}
			httpReq.RequestURI = url
		{{if HasUpload . }}
		{{else if HasInput . }}
			httpReq.Header.Set("Content-type", "application/json")
		{{end}}
		{{if HasOutput . }}
			httpReq.Header.Set("Accept", "application/json")
		{{end}}
		for k, v := range headers {
			httpReq.Header.Set(k, v)
		}
		setCookieHook(httpReq, headers)

		// record request-part of test-case
		tc.WithRequest("GET", url, httpReq.Header, []byte{})
	}

	// call server
	httpResp := httptest.NewRecorder()
	{
	    // invoke business logic on remote service
    	webservice := NewRest{{ToFirstUpper $serviceName}}()
    	webservice.HTTPHandler().ServeHTTP(httpResp, httpReq)
	}

	{
		// record response-part of test-case
		tc.WithResponse(httpResp.Code, httpResp.Header(), httpResp.Body.Bytes())
	}

	// handle response
	{{if IsRestOperationJSON . }}
		{{if HasOutput . }}
				if httpResp.Code != http.StatusOK {
					// return error response
					var errorResp errorh.Error
					dec := json.NewDecoder(httpResp.Body)
					err = dec.Decode(&errorResp)
					if err != nil {
						t.Fatalf("Error unmarshalling response: %s", err )
					}
					return httpResp.Code, nil, &errorResp, nil
				}

				// return success response
				resp := {{GetOutputArgDeclaration . }}
				dec := json.NewDecoder(httpResp.Body)
				err = dec.Decode({{GetOutputArgName . }})
				if err != nil {
					return httpResp.Code, nil, nil, err
				}
				return httpResp.Code, resp, nil, nil
		{{else}}
			return httpResp.Code, nil, nil
		{{end}}
	{{else}}
		return httpResp, nil
	{{end}}
}
    {{end}}
{{end}}

func defaultBeforeAll() {
    mytime.SetMockNow()
}

func defaultAfterAll() {
    mytime.SetDefaultNow()
}

func fetchEvents(c context.Context) []string {
	found := []string{}
	eventStore.Mocked().IterateAll(c, rest.Credentials{}, func(e envelope.Envelope) error {
		found = append(found, fmt.Sprintf("%s.%s", e.AggregateName, e.EventTypeName))
		return nil
	})
	return found
}
`
