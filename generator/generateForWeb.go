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
	{
		target := fmt.Sprintf("%s/httpHandlers.go", targetDir)

		data := Structs{
			PackageName: packageName,
			Structs:     structs,
		}
		err = generateFileFromTemplate(data, "handlers", target)
		if err != nil {
			log.Fatalf("Error generating wrappers for structs (%s)", err)
			return err
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

  "github.com/gorilla/mux"
  "github.com/MarcGrol/microgen/lib/myerrors"
)

{{range .Structs}}
	{{if .IsRestService}}
	func (ts *{{.Name}}) HandleHttp() http.Handler {
		router := mux.NewRouter().StrictSlash(true)
		{{range .Operations}}
			{{if .IsRestOperation}}
				router.HandleFunc("{{.GetRestOperationPath}}", {{.Name}}(ts)).Methods("{{.GetRestOperationMethod}}")
			{{end}}
		{{end}}
		return router
	}
	{{end}}
{{end}}

{{range $idx, $str := .Structs}}
	{{if .IsRestService}}
		{{range $idxOper, $oper := .Operations}}
			{{if .IsRestOperation}}
				 func {{$oper.Name}}( service *{{$str.Name}} ) http.HandlerFunc {
						return func(w http.ResponseWriter, r *http.Request) {
							var err error

							pathParams := mux.Vars(r)
							log.Printf("pathParams:%+v", pathParams)

							// extract url-params
							{{range .InputArgs}}
								{{if .IsPrimitive}}
								{{.Name}}, exists := pathParams["{{.Name}}"]
								if !exists {
									handleError(myerrors.NewInvalidInputError(fmt.Errorf("Missing path param '{{.Name}}'")), w)
									return
								}
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
	{{end}}
{{end}}
`
