// Generated automatically: do not edit manually

package testData

import (
	"net/http"
	"net/http/httptest"
)

func doitTestHelper(url string) (int, error) {

	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {

		return 0, err

	}

	webservice := MyService{}
	webservice.HttpHandler().ServeHTTP(recorder, req)

	return recorder.Code, nil

}
