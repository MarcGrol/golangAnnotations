package rest

const testServiceTemplate = `package {{.PackageName}}

// Generated automatically by golangAnnotations: do not edit manually

// HTTPTestHandlerWithRouter registers endpoint in existing router
func HTTPTestHandlerWithRouter(router *mux.Router, jsonTestResults string) *mux.Router {
	subRouter := router.PathPrefix("{{GetRestServicePath . }}").Subrouter()

	suite := libtest.HTTPTestSuite{}
	err := json.Unmarshal([]byte(jsonTestResults), &suite)
	if err == nil {
		subRouter.HandleFunc("/logs.md", writeTestLogsAsMarkdown(suite)).Methods("GET")
	}

	return router
}

func writeTestLogsAsMarkdown(suite libtest.HTTPTestSuite) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/markdown; charset=UTF-8")
		fmt.Fprintf(w, "%s", suite.AsMarkdown())
	}
}
`
