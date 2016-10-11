
// Generated automatically by golangAnnotations: do not edit manually

package testData

import (
"github.com/Duxxie/platform/backend/lib/testcase"
"github.com/gorilla/mux"
)

// HTTPTestHandlerWithRouter registers endpoint in existing router
func HTTPTestHandlerWithRouter(router *mux.Router, results testcase.TestSuiteDescriptor) *mux.Router {
	subRouter := router.PathPrefix("/api").Subrouter()

	subRouter.HandleFunc("/logs.md", testcase.WriteTestLogsAsMarkdown(results)).Methods("GET")

	return router
}

