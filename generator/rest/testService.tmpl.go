package rest

const testServiceTemplate = `// Generated automatically by golangAnnotations: do not edit manually

package {{.PackageName}}

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var testResults = ""

// HTTPTestHandlerWithRouter registers endpoint in existing router
func HTTPTestHandlerWithRouter(router *mux.Router) *mux.Router {
	subRouter := router.PathPrefix("{{GetRestServicePath . }}").Subrouter()

	subRouter.HandleFunc("/logs.md", writeTestLogsAsMarkdown()).Methods("GET")

	return router
}

func writeTestLogsAsMarkdown() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/markdown; charset=UTF-8")
		fmt.Fprintf(w, "%s", testResults)
	}
}
`
