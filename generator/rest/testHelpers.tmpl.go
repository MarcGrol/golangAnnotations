package rest

const testHelpersTemplate = `// +build !appengine

// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

import (
    "golang.org/x/net/context"
)

var (
    setCookieHook = func(r *http.Request, headers map[string]string) {}
	testSet= NewHTTPTestSet("{{.PackageName}}")
)

func TestMain(m *testing.M) {
     beforeAll()

     code := m.Run()

     afterAll()

     testSet.WriteToFile()

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

{{ $serviceName := .Name }}

{{range .Operations}}

{{if IsRestOperation . }}
func {{.Name}}TestHelper(t *testing.T, c context.Context, tc *HTTPTestCase, url string {{if IsRestOperationForm . }}, form url.Values{{else if HasInput . }}, input {{GetInputArgType . }} {{end}} )  ({{if IsRestOperationJSON . }}int {{if HasOutput . }},{{GetOutputArgType . }}{{end}},*errorh.Error{{else}}*httptest.ResponseRecorder{{end}},error) {
    return {{.Name}}TestHelperWithHeaders( t, c, tc, url {{if IsRestOperationForm . }}, form{{else if HasInput . }}, input {{end}}, map[string]string{} )
}

func {{.Name}}TestHelperWithHeaders(t *testing.T, c context.Context,  tc *HTTPTestCase, url string {{if IsRestOperationForm . }}, form url.Values{{else if HasInput . }}, input {{GetInputArgType . }} {{end}}, headers map[string]string)  ({{if IsRestOperationJSON . }}int {{if HasOutput . }},{{GetOutputArgType . }}{{end}},*errorh.Error{{else}}*httptest.ResponseRecorder{{end}},error) {
	// collect test-case info
    tc.WithOperationName("{{.Name}}").
       WithPreConditions([]string{"eventA","eventB"})
	defer func() {
        tc.WithPreConditions([]string{"eventC"})
        testSet.Add(tc)
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

    // invoke remote service
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
	{{else}}  // else HasOutput
		return httpResp.Code, nil
	{{end}} // end IsRestOperationJSON
}
    {{end}} // end IsRestOperation
{{end}} // end range
`
