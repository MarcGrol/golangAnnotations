package rest

const httpClientTemplate = `// +build !appengine

// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"
)

{{ $serviceName := .Name }}

var debug = false

type HTTPClient struct {
	hostName string
}

func NewHTTPClient(host string) *HTTPClient {
	return &HTTPClient{
		hostName: host,
	}
}

{{range .Operations -}}

{{if IsRestOperation . -}}
	{{if IsRestOperationJSON . -}}

// {{ToFirstUpper .Name}} can be used by external clients to interact with the system
func (c *HTTPClient) {{ToFirstUpper .Name}}(ctx context.Context, url string{{if HasInput . }}, input {{GetInputArgType . }}{{end}}, cookie *http.Cookie, requestUID string, timeout time.Duration) (int{{if HasOutput . }}, {{GetOutputArgType . }}{{end}}, *errorh.Error, error) {

	{{if HasInput . -}}
	requestBody, _ := json.Marshal(input)
	req, err := http.NewRequest("{{GetRestOperationMethod . }}", c.hostName+url, strings.NewReader(string(requestBody)))
	{{else -}}
	req, err := http.NewRequest("{{GetRestOperationMethod . }}", c.hostName+url, nil)
	{{end -}}
	if err != nil {
		{{if HasOutput . -}}
			return 0, nil, nil, err
		{{else -}}
			return 0, nil, err
		{{end -}}
	}
	if cookie != nil {
		req.AddCookie(cookie)
	}
	{{if HasInput . -}}
		req.Header.Set("Content-type", "application/json")
	{{end -}}
	{{if HasOutput . -}}
	req.Header.Set("Accept", "application/json")
	{{end -}}
	req.Header.Set("X-CSRF-Token", "true")

	if debug {
		dump, err := httputil.DumpRequest(req, true)
		if err == nil {
			mylog.New().Debug(ctx, "HTTP request-payload:\n %s", dump)
		}
	}

	cl := http.Client{}
	cl.Timeout = timeout
	res, err := cl.Do(req)
	if err != nil {
		{{if HasOutput . -}}
		return -1, nil, nil, err
		{{else -}}
		return -1, nil, nil
	{{end -}}
	}
	defer res.Body.Close()

	if debug {
		respDump, err := httputil.DumpResponse(res, true)
		if err == nil {
			mylog.New().Debug(ctx, "HTTP response-payload:\n%s", string(respDump))
		}
	}

	{{if HasOutput . -}}
	if res.StatusCode >= http.StatusMultipleChoices {
		// return error response
		var errorResp errorh.Error
		dec := json.NewDecoder(res.Body)
		err = dec.Decode(&errorResp)
		if err != nil {
			return res.StatusCode, nil, nil, err
		}
		return res.StatusCode, nil, &errorResp, nil
	}

	// return success response
	resp := {{GetOutputArgDeclaration . }}
	dec := json.NewDecoder(res.Body)
	err = dec.Decode({{GetOutputArgName . }})
	if err != nil {
		return res.StatusCode, nil, nil, err
	}
	return res.StatusCode, resp, nil, nil

	{{else -}}
	return res.StatusCode, nil, nil
	{{end -}}
}
		{{end -}}
	{{end -}}
{{end -}}
`
