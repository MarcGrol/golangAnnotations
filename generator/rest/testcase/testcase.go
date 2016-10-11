package testcase

import (
	"fmt"
	"net/http"
	"strings"
)

type TestSuiteDescriptor struct {
	TestCases []TestCaseDescriptor
}

type TestCaseDescriptor struct {
	Name        string
	Description string
	Operation   string
	Request     RequestDescriptor
	Response    ResponseDescriptor
}

type RequestDescriptor struct {
	Method  string
	Url     string
	Headers []string
	Body    string
}

type ResponseDescriptor struct {
	Status  int
	Headers []string
	Body    string
}

// WriteTestLogsAsMarkdown returns the test reports of a package as markdown
func WriteTestLogsAsMarkdown(results TestSuiteDescriptor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/markdown")
		fmt.Fprintf(w, "# HTTP exchange for all testcases in package\n\n")
		for _, tc := range results.TestCases {
			fmt.Fprintf(w, "## %s\n", tc.Name)
			fmt.Fprintf(w, "%s\n", tc.Description)
			fmt.Fprintf(w, "### Operation %s\n", tc.Operation)

			fmt.Fprintf(w, "\n### http-request:\n")
			fmt.Fprintf(w, "    %s %s%s HTTP/1.1\n", tc.Request.Method, "/api", tc.Request.Url)
			for _, h := range tc.Request.Headers {
				fmt.Fprintf(w, "    %s\n", h)
			}
			fmt.Fprintf(w, "\n    %s\n", strings.Replace(string(tc.Request.Body), "\n", "\n    ", -1))

			fmt.Fprintf(w, "\n### http-response:\n\n")
			fmt.Fprintf(w, "    %d\n", tc.Response.Status)
			for _, h := range tc.Response.Headers {
				fmt.Fprintf(w, "    %s\n", h)
			}
			fmt.Fprintf(w, "\n    %s\n\n", strings.Replace(tc.Response.Body, "\n", "\n    ", -1))
			fmt.Fprintf(w, "---\n\n")

		}

	}
}
