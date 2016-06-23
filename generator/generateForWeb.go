package generator

import (
	"fmt"
	"log"

	"github.com/MarcGrol/astTools/model"
)

func GenerateForWeb(inputDir string, structs []model.Struct) error {
	packageName, err := getPackageName(structs)
	if err != nil {
		return err
	}
	targetDir, err := determineTargetPath(inputDir, packageName)
	if err != nil {
		return err
	}
	for _, service := range structs {
		if service.IsRestService() {
			{
				target := fmt.Sprintf("%s/http%s.go", targetDir, service.Name)
				err = generateFileFromTemplate(service, "handlers", target)
				if err != nil {
					log.Fatalf("Error generating handlers for service %s: %s", service.Name, err)
					return err
				}
			}
			{
				target := fmt.Sprintf("%s/http%sHelpers_test.go", targetDir, service.Name)
				err = generateFileFromTemplate(service, "helpers", target)
				if err != nil {
					log.Fatalf("Error generating helpers for service %s: %s", service.Name, err)
					return err
				}
			}

		}
	}
	return nil
}

var handlersTemplate string = `
// Generated automatically: do not edit manually

package {{.PackageName}}

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/MarcGrol/microgen/lib/myerrors"
	"github.com/gorilla/mux"
)

{{ $structName := .Name }}

func (ts *{{.Name}}) HttpHandler() http.Handler {
	servicePrefix := "{{.GetRestServicePath}}"
	router := mux.NewRouter().StrictSlash(true)
	{{range .Operations}}
		{{if .IsRestOperation}}
			router.HandleFunc( servicePrefix + "{{.GetRestOperationPath}}", {{.Name}}(ts)).Methods("{{.GetRestOperationMethod}}")
		{{end}}
	{{end}}
	return router
}

{{range $idxOper, $oper := .Operations}}
{{if .IsRestOperation}}
func {{$oper.Name}}( service *{{$structName}} ) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		pathParams := mux.Vars(r)
		log.Printf("pathParams:%+v", pathParams)

		// extract url-params
		{{range .InputArgs}}
			{{if .IsPrimitive}}
				{{if .IsNumber}}
					{{.Name}}String, exists := pathParams["{{.Name}}"]
					if !exists {
						handleError(myerrors.NewInvalidInputError(fmt.Errorf("Missing path param '{{.Name}}'")), w)
						return
					}
					{{.Name}}, err := strconv.Atoi({{.Name}}String)
					if err != nil {
						handleError(myerrors.NewInvalidInputError(fmt.Errorf("Invalid path param '{{.Name}}'")), w)
						return
					}
				{{else}}
					{{.Name}}, exists := pathParams["{{.Name}}"]
					if !exists {
						handleError(myerrors.NewInvalidInputError(fmt.Errorf("Missing path param '{{.Name}}'")), w)
						return
					}
				{{end}}
			{{end}}
		{{end}}

		{{if .HasInput }}
			// read abd parse request body
			var {{.GetInputArgName}} {{.GetInputArgType}}
			err = json.NewDecoder(r.Body).Decode( &{{.GetInputArgName}} )
			if err != nil {
				handleError(myerrors.NewInvalidInputError(fmt.Errorf("Error decoding request payload:%s", err)), w)
				return
			}
		{{end}}

		// call business logic
		{{if .HasOutput }}
			result, err := service.{{$oper.Name}}({{.GetInputParamString}})
		{{else}}
			err = service.{{$oper.Name}}({{.GetInputParamString}})
		{{end}}
		if err != nil {
			handleError(err, w)
			return
		}

		// write response body
		{{if .HasOutput }}
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(result)
			if err != nil {
				log.Printf("Error encoding response payload %+v", err)
			}
		{{else}}
			w.WriteHeader(http.StatusNoContent)
		{{end}}
      }
 }
{{end}}
{{end}}


func handleError(err error, w http.ResponseWriter) {
	errorBody := struct {
		ErrorMessage string
	}{
		err.Error(),
	}
	blob, err := json.Marshal(errorBody)
	if err != nil {
		log.Printf("Error marshalling error response payload %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(determineHttpCode(err))
	w.Header().Set("Content-Type", "application/json")
	w.Write(blob)
}

func determineHttpCode(err error) int {
	if myerrors.IsNotFoundError(err) {
		return http.StatusNotFound
	} else if myerrors.IsInternalError(err) {
		return http.StatusInternalServerError
	} else if myerrors.IsInvalidInputError(err) {
		return http.StatusBadRequest
	} else if myerrors.IsNotAuthorizedError(err) {
		return http.StatusForbidden
	} else {
		return http.StatusInternalServerError
	}
}

`

var helpersTemplate string = `
// Generated automatically: do not edit manually

package {{.PackageName}}

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
)

{{ $structName := .Name }}

{{range .Operations}}

{{if .IsRestOperation}}
func {{.Name}}TestHelper(url string {{if .HasInput}}, input {{.GetInputArgType}} {{end}} )  (int {{if .HasOutput }},*{{.GetOutputArgType}}{{end}},error) {

	recorder := httptest.NewRecorder()

	{{if .HasInput}}
		requestBody, _ := json.Marshal(input)
		req, err := http.NewRequest("{{.GetRestOperationMethod}}", url, strings.NewReader(string(requestBody)))
	{{else}}
		req, err := http.NewRequest("{{.GetRestOperationMethod}}", url, nil)
	{{end}}
	if err != nil {
		{{if .HasOutput }}
			return 0, nil, err
		{{else}}
			return 0,  err
		{{end}}
	}
	{{if .HasOutput }}
		req.Header.Set("Accept", "application/json")
	{{end}}

	webservice := {{$structName}}{}
	webservice.HttpHandler().ServeHTTP(recorder, req)

	{{if .HasOutput }}
		var resp {{.GetOutputArgType}}
		dec := json.NewDecoder(recorder.Body)
		err = dec.Decode(&resp)
		if err != nil {
			return recorder.Code, nil, err
		}
		return recorder.Code, &resp, nil
	{{else}}
		return recorder.Code, nil
	{{end}}
}
{{end}}
{{end}}
`
