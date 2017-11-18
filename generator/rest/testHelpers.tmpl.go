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

func defaultBeforeAll() {
    mytime.SetMockNow()
}

func defaultAfterAll() {
    mytime.SetDefaultNow()
}

func TestMain(m *testing.M) {
	beforeAll()
	code := m.Run()
	afterAll()
	testSuite.WriteToMarkdownGoVarFile()
	//testSuite.WriteToJsonFile()
	os.Exit(code)
}

{{ $serviceName := .Name }}

{{range .Operations}}

{{if IsRestOperation . }}
func {{.Name}}TestHelper(t *testing.T, c context.Context, tc *libtest.HTTPTestCase, url string {{if IsRestOperationForm . }}, form url.Values{{else if HasInput . }}, input {{GetInputArgType . }} {{end}} )  ({{if IsRestOperationJSON . }}int {{if HasOutput . }},{{GetOutputArgType . }}{{end}},*errorh.Error{{else}}*httptest.ResponseRecorder{{end}},error) {
    return {{.Name}}TestHelperWithHeaders( t, c, tc, url {{if IsRestOperationForm . }}, form{{else if HasInput . }}, input {{end}}, map[string]string{} )
}

func {{.Name}}TestHelperWithHeaders(t *testing.T, c context.Context,  tc *libtest.HTTPTestCase, url string {{if IsRestOperationForm . }}, form url.Values{{else if HasInput . }}, input {{GetInputArgType . }} {{end}}, headers map[string]string)  ({{if IsRestOperationJSON . }}int {{if HasOutput . }},{{GetOutputArgType . }}{{end}},*errorh.Error{{else}}*httptest.ResponseRecorder{{end}},error) {
	// collect test-case info
    tc.WithOperationName("{{.Name}}").
       WithPreConditions([]string{"eventA","eventB"}) // TODO collect events from store
	defer func() {
        tc.WithPostConditions([]string{"eventC"}) // TODO collect delta events from store
        testSuite.Add(tc)
	}()

    // create http-request
    {{if HasUpload . }}
        {{.Name}}SetUpload(input)
        httpReq, err := http.NewRequest("{{GetRestOperationMethod . }}", url, nil)
    {{else if IsRestOperationForm . }}
        httpReq, err := http.NewRequest("{{GetRestOperationMethod . }}", url, strings.NewReader(form.Encode()))
        httpReq.Header.Add("Content-Type", "application/x-www-form-urlencoded")
    {{else if HasInput . }}
        rb, _ := json.Marshal(input)
        // indent for readability
        var requestBody bytes.Buffer
        json.Indent(&requestBody, rb, "", "\t")
        httpReq, err := http.NewRequest("{{GetRestOperationMethod . }}", url, strings.NewReader(requestBody.String()))
    {{else}}
        httpReq, err := http.NewRequest("{{GetRestOperationMethod . }}", url, nil)
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

    // collect test-case info
    tc.WithRequest("GET", url, httpReq.Header, []byte{})

    //
    // invoke business logic as remote service
	//
    httpResp := httptest.NewRecorder()
    webservice := NewRest{{ToFirstUpper $serviceName}}()
    webservice.HTTPHandler().ServeHTTP(httpResp, httpReq)

	// collect response related test-case info
	tc.WithResponse(httpResp.Code, httpResp.Header(), httpResp.Body.Bytes())

	{{if IsRestOperationJSON . }}
		{{if HasOutput . }}
				if httpResp.Code != http.StatusOK {
					// return error response
					var errorResp errorh.Error
					dec := json.NewDecoder(httpResp.Body)
					err = dec.Decode(&errorResp)
					if err != nil {
						return httpResp.Code, nil, nil, err
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
`
