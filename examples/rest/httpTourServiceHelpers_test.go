// +build !appengine

// Generated automatically by golangAnnotations: do not edit manually

package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/Duxxie/platform/backend/lib/mytime"
	"github.com/MarcGrol/golangAnnotations/generator/rest/errorh"
)

var logFp *os.File

func openfile(filename string) *os.File {
	fp, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Error opening rest-dump-file %s: %s", filename, err.Error())
	}
	return fp
}

func TestMain(m *testing.M) {

	dirname := "restTestLog"
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		os.Mkdir(dirname, os.ModePerm)
	}
	logFp = openfile(dirname + "/testResults.go")
	defer func() {
		logFp.Close()
	}()
	fmt.Fprintf(logFp, "package %s\n\n", dirname)
	fmt.Fprintf(logFp, "// Generated automatically based on running of api-tests\n\n")
	fmt.Fprintf(logFp, "import (\n")
	fmt.Fprintf(logFp, "\"github.com/MarcGrol/golangAnnotations/generator/rest/testcase\"\n")
	fmt.Fprintf(logFp, ")\n")

	fmt.Fprintf(logFp, "var TestResults = testcase.TestSuiteDescriptor {\n")
	fmt.Fprintf(logFp, "\tTestCases: []testcase.TestCaseDescriptor{\n")

	beforeAll()

	code := m.Run()

	afterAll()

	fmt.Fprintf(logFp, "},\n")
	fmt.Fprintf(logFp, "}\n")

	os.Exit(code)
}

func beforeAll() {
	mytime.SetMockNow()
}

func afterAll() {
	mytime.SetDefaultNow()
}

func testCase(name string, description string) {
	fmt.Fprintf(logFp, "\t\ttestcase.TestCaseDescriptor{\n")
	fmt.Fprintf(logFp, "\t\tName:\"%s\",\n", name)
	fmt.Fprintf(logFp, "\t\tDescription:\"%s\",\n", description)
}

func testCaseDone() {
	fmt.Fprintf(logFp, "},\n")
}

func getTourOnUidTestHelper(url string) (int, *Tour, *errorh.Error, error) {
	return getTourOnUidTestHelperWithHeaders(url, map[string]string{})
}

func getTourOnUidTestHelperWithHeaders(url string, headers map[string]string) (int, *Tour, *errorh.Error, error) {

	fmt.Fprintf(logFp, "\t\tOperation:\"%s\",\n", "getTourOnUid")
	defer func() {
		fmt.Fprintf(logFp, "\t},\n")
	}()

	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {

		return 0, nil, nil, err

	}
	req.RequestURI = url

	req.Header.Set("Accept", "application/json")

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	headersToBeSorted := []string{}
	for key, values := range req.Header {
		for _, value := range values {
			headersToBeSorted = append(headersToBeSorted, fmt.Sprintf("%s:%s", key, value))
		}
	}
	sort.Strings(headersToBeSorted)

	fmt.Fprintf(logFp, "\tRequest: testcase.RequestDescriptor{\n")
	fmt.Fprintf(logFp, "\tMethod:\"%s\",\n", "GET")
	fmt.Fprintf(logFp, "\tUrl:\"%s\",\n", url)
	fmt.Fprintf(logFp, "\tHeaders: []string{\n")
	for _, h := range headersToBeSorted {
		fmt.Fprintf(logFp, "\"%s\",\n", h)
	}
	fmt.Fprintf(logFp, "\t},\n")

	fmt.Fprintf(logFp, "},\n")

	// dump readable request
	//payload, err := httputil.DumpRequest(req, true)

	fmt.Fprintf(logFp, "\tResponse:testcase.ResponseDescriptor{\n")
	defer func() {
		fmt.Fprintf(logFp, "\t},\n")
	}()

	webservice := TourService{}
	webservice.HTTPHandler().ServeHTTP(recorder, req)

	// dump readable response
	var responseBody bytes.Buffer
	json.Indent(&responseBody, recorder.Body.Bytes(), "", "\t")

	fmt.Fprintf(logFp, "\tStatus:%d,\n", recorder.Code)

	headersToBeSorted = []string{}
	for key, values := range recorder.Header() {
		for _, value := range values {
			headersToBeSorted = append(headersToBeSorted, fmt.Sprintf("%s:%s", key, value))
		}
	}
	sort.Strings(headersToBeSorted)

	fmt.Fprintf(logFp, "\tHeaders:[]string{\n")
	for _, h := range headersToBeSorted {
		fmt.Fprintf(logFp, "\"%s\",\n", h)
	}
	fmt.Fprintf(logFp, "\t},\n")
	fmt.Fprintf(logFp, "\tBody:\n`%s`,\n", responseBody.String())

	if recorder.Code != http.StatusOK {
		// return error response
		var errorResp errorh.Error
		dec := json.NewDecoder(recorder.Body)
		err = dec.Decode(&errorResp)
		if err != nil {
			return recorder.Code, nil, nil, err
		}
		return recorder.Code, nil, &errorResp, nil
	}

	// return success response
	resp := &Tour{}
	dec := json.NewDecoder(recorder.Body)
	err = dec.Decode(resp)
	if err != nil {
		return recorder.Code, nil, nil, err
	}
	return recorder.Code, resp, nil, nil

}

func createEtappeTestHelper(url string, input Etappe) (int, *Etappe, *errorh.Error, error) {
	return createEtappeTestHelperWithHeaders(url, input, map[string]string{})
}

func createEtappeTestHelperWithHeaders(url string, input Etappe, headers map[string]string) (int, *Etappe, *errorh.Error, error) {

	fmt.Fprintf(logFp, "\t\tOperation:\"%s\",\n", "createEtappe")
	defer func() {
		fmt.Fprintf(logFp, "\t},\n")
	}()

	recorder := httptest.NewRecorder()

	rb, _ := json.Marshal(input)
	// indent for readability
	var requestBody bytes.Buffer
	json.Indent(&requestBody, rb, "", "\t")

	req, err := http.NewRequest("POST", url, strings.NewReader(requestBody.String()))

	if err != nil {

		return 0, nil, nil, err

	}
	req.RequestURI = url

	req.Header.Set("Accept", "application/json")

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	headersToBeSorted := []string{}
	for key, values := range req.Header {
		for _, value := range values {
			headersToBeSorted = append(headersToBeSorted, fmt.Sprintf("%s:%s", key, value))
		}
	}
	sort.Strings(headersToBeSorted)

	fmt.Fprintf(logFp, "\tRequest: testcase.RequestDescriptor{\n")
	fmt.Fprintf(logFp, "\tMethod:\"%s\",\n", "POST")
	fmt.Fprintf(logFp, "\tUrl:\"%s\",\n", url)
	fmt.Fprintf(logFp, "\tHeaders: []string{\n")
	for _, h := range headersToBeSorted {
		fmt.Fprintf(logFp, "\"%s\",\n", h)
	}
	fmt.Fprintf(logFp, "\t},\n")

	fmt.Fprintf(logFp, "\tBody:\n")
	fmt.Fprintf(logFp, "`%s`", requestBody.String())
	fmt.Fprintf(logFp, ",\n")

	fmt.Fprintf(logFp, "},\n")

	// dump readable request
	//payload, err := httputil.DumpRequest(req, true)

	fmt.Fprintf(logFp, "\tResponse:testcase.ResponseDescriptor{\n")
	defer func() {
		fmt.Fprintf(logFp, "\t},\n")
	}()

	webservice := TourService{}
	webservice.HTTPHandler().ServeHTTP(recorder, req)

	// dump readable response
	var responseBody bytes.Buffer
	json.Indent(&responseBody, recorder.Body.Bytes(), "", "\t")

	fmt.Fprintf(logFp, "\tStatus:%d,\n", recorder.Code)

	headersToBeSorted = []string{}
	for key, values := range recorder.Header() {
		for _, value := range values {
			headersToBeSorted = append(headersToBeSorted, fmt.Sprintf("%s:%s", key, value))
		}
	}
	sort.Strings(headersToBeSorted)

	fmt.Fprintf(logFp, "\tHeaders:[]string{\n")
	for _, h := range headersToBeSorted {
		fmt.Fprintf(logFp, "\"%s\",\n", h)
	}
	fmt.Fprintf(logFp, "\t},\n")
	fmt.Fprintf(logFp, "\tBody:\n`%s`,\n", responseBody.String())

	if recorder.Code != http.StatusOK {
		// return error response
		var errorResp errorh.Error
		dec := json.NewDecoder(recorder.Body)
		err = dec.Decode(&errorResp)
		if err != nil {
			return recorder.Code, nil, nil, err
		}
		return recorder.Code, nil, &errorResp, nil
	}

	// return success response
	resp := &Etappe{}
	dec := json.NewDecoder(recorder.Body)
	err = dec.Decode(resp)
	if err != nil {
		return recorder.Code, nil, nil, err
	}
	return recorder.Code, resp, nil, nil

}

func addEtappeResultsTestHelper(url string, input EtappeResult) (int, *errorh.Error, error) {
	return addEtappeResultsTestHelperWithHeaders(url, input, map[string]string{})
}

func addEtappeResultsTestHelperWithHeaders(url string, input EtappeResult, headers map[string]string) (int, *errorh.Error, error) {

	fmt.Fprintf(logFp, "\t\tOperation:\"%s\",\n", "addEtappeResults")
	defer func() {
		fmt.Fprintf(logFp, "\t},\n")
	}()

	recorder := httptest.NewRecorder()

	rb, _ := json.Marshal(input)
	// indent for readability
	var requestBody bytes.Buffer
	json.Indent(&requestBody, rb, "", "\t")

	req, err := http.NewRequest("PUT", url, strings.NewReader(requestBody.String()))

	if err != nil {

		return 0, nil, err

	}
	req.RequestURI = url

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	headersToBeSorted := []string{}
	for key, values := range req.Header {
		for _, value := range values {
			headersToBeSorted = append(headersToBeSorted, fmt.Sprintf("%s:%s", key, value))
		}
	}
	sort.Strings(headersToBeSorted)

	fmt.Fprintf(logFp, "\tRequest: testcase.RequestDescriptor{\n")
	fmt.Fprintf(logFp, "\tMethod:\"%s\",\n", "PUT")
	fmt.Fprintf(logFp, "\tUrl:\"%s\",\n", url)
	fmt.Fprintf(logFp, "\tHeaders: []string{\n")
	for _, h := range headersToBeSorted {
		fmt.Fprintf(logFp, "\"%s\",\n", h)
	}
	fmt.Fprintf(logFp, "\t},\n")

	fmt.Fprintf(logFp, "\tBody:\n")
	fmt.Fprintf(logFp, "`%s`", requestBody.String())
	fmt.Fprintf(logFp, ",\n")

	fmt.Fprintf(logFp, "},\n")

	// dump readable request
	//payload, err := httputil.DumpRequest(req, true)

	fmt.Fprintf(logFp, "\tResponse:testcase.ResponseDescriptor{\n")
	defer func() {
		fmt.Fprintf(logFp, "\t},\n")
	}()

	webservice := TourService{}
	webservice.HTTPHandler().ServeHTTP(recorder, req)

	// dump readable response
	var responseBody bytes.Buffer
	json.Indent(&responseBody, recorder.Body.Bytes(), "", "\t")

	fmt.Fprintf(logFp, "\tStatus:%d,\n", recorder.Code)

	headersToBeSorted = []string{}
	for key, values := range recorder.Header() {
		for _, value := range values {
			headersToBeSorted = append(headersToBeSorted, fmt.Sprintf("%s:%s", key, value))
		}
	}
	sort.Strings(headersToBeSorted)

	fmt.Fprintf(logFp, "\tHeaders:[]string{\n")
	for _, h := range headersToBeSorted {
		fmt.Fprintf(logFp, "\"%s\",\n", h)
	}
	fmt.Fprintf(logFp, "\t},\n")
	fmt.Fprintf(logFp, "\tBody:\n`%s`,\n", responseBody.String())

	return recorder.Code, nil, nil

}

func createCyclistTestHelper(url string, input Cyclist) (int, *Cyclist, *errorh.Error, error) {
	return createCyclistTestHelperWithHeaders(url, input, map[string]string{})
}

func createCyclistTestHelperWithHeaders(url string, input Cyclist, headers map[string]string) (int, *Cyclist, *errorh.Error, error) {

	fmt.Fprintf(logFp, "\t\tOperation:\"%s\",\n", "createCyclist")
	defer func() {
		fmt.Fprintf(logFp, "\t},\n")
	}()

	recorder := httptest.NewRecorder()

	rb, _ := json.Marshal(input)
	// indent for readability
	var requestBody bytes.Buffer
	json.Indent(&requestBody, rb, "", "\t")

	req, err := http.NewRequest("POST", url, strings.NewReader(requestBody.String()))

	if err != nil {

		return 0, nil, nil, err

	}
	req.RequestURI = url

	req.Header.Set("Accept", "application/json")

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	headersToBeSorted := []string{}
	for key, values := range req.Header {
		for _, value := range values {
			headersToBeSorted = append(headersToBeSorted, fmt.Sprintf("%s:%s", key, value))
		}
	}
	sort.Strings(headersToBeSorted)

	fmt.Fprintf(logFp, "\tRequest: testcase.RequestDescriptor{\n")
	fmt.Fprintf(logFp, "\tMethod:\"%s\",\n", "POST")
	fmt.Fprintf(logFp, "\tUrl:\"%s\",\n", url)
	fmt.Fprintf(logFp, "\tHeaders: []string{\n")
	for _, h := range headersToBeSorted {
		fmt.Fprintf(logFp, "\"%s\",\n", h)
	}
	fmt.Fprintf(logFp, "\t},\n")

	fmt.Fprintf(logFp, "\tBody:\n")
	fmt.Fprintf(logFp, "`%s`", requestBody.String())
	fmt.Fprintf(logFp, ",\n")

	fmt.Fprintf(logFp, "},\n")

	// dump readable request
	//payload, err := httputil.DumpRequest(req, true)

	fmt.Fprintf(logFp, "\tResponse:testcase.ResponseDescriptor{\n")
	defer func() {
		fmt.Fprintf(logFp, "\t},\n")
	}()

	webservice := TourService{}
	webservice.HTTPHandler().ServeHTTP(recorder, req)

	// dump readable response
	var responseBody bytes.Buffer
	json.Indent(&responseBody, recorder.Body.Bytes(), "", "\t")

	fmt.Fprintf(logFp, "\tStatus:%d,\n", recorder.Code)

	headersToBeSorted = []string{}
	for key, values := range recorder.Header() {
		for _, value := range values {
			headersToBeSorted = append(headersToBeSorted, fmt.Sprintf("%s:%s", key, value))
		}
	}
	sort.Strings(headersToBeSorted)

	fmt.Fprintf(logFp, "\tHeaders:[]string{\n")
	for _, h := range headersToBeSorted {
		fmt.Fprintf(logFp, "\"%s\",\n", h)
	}
	fmt.Fprintf(logFp, "\t},\n")
	fmt.Fprintf(logFp, "\tBody:\n`%s`,\n", responseBody.String())

	if recorder.Code != http.StatusOK {
		// return error response
		var errorResp errorh.Error
		dec := json.NewDecoder(recorder.Body)
		err = dec.Decode(&errorResp)
		if err != nil {
			return recorder.Code, nil, nil, err
		}
		return recorder.Code, nil, &errorResp, nil
	}

	// return success response
	resp := &Cyclist{}
	dec := json.NewDecoder(recorder.Body)
	err = dec.Decode(resp)
	if err != nil {
		return recorder.Code, nil, nil, err
	}
	return recorder.Code, resp, nil, nil

}

func markCyclistAbondonedTestHelper(url string) (int, *errorh.Error, error) {
	return markCyclistAbondonedTestHelperWithHeaders(url, map[string]string{})
}

func markCyclistAbondonedTestHelperWithHeaders(url string, headers map[string]string) (int, *errorh.Error, error) {

	fmt.Fprintf(logFp, "\t\tOperation:\"%s\",\n", "markCyclistAbondoned")
	defer func() {
		fmt.Fprintf(logFp, "\t},\n")
	}()

	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("DELETE", url, nil)

	if err != nil {

		return 0, nil, err

	}
	req.RequestURI = url

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	headersToBeSorted := []string{}
	for key, values := range req.Header {
		for _, value := range values {
			headersToBeSorted = append(headersToBeSorted, fmt.Sprintf("%s:%s", key, value))
		}
	}
	sort.Strings(headersToBeSorted)

	fmt.Fprintf(logFp, "\tRequest: testcase.RequestDescriptor{\n")
	fmt.Fprintf(logFp, "\tMethod:\"%s\",\n", "DELETE")
	fmt.Fprintf(logFp, "\tUrl:\"%s\",\n", url)
	fmt.Fprintf(logFp, "\tHeaders: []string{\n")
	for _, h := range headersToBeSorted {
		fmt.Fprintf(logFp, "\"%s\",\n", h)
	}
	fmt.Fprintf(logFp, "\t},\n")

	fmt.Fprintf(logFp, "},\n")

	// dump readable request
	//payload, err := httputil.DumpRequest(req, true)

	fmt.Fprintf(logFp, "\tResponse:testcase.ResponseDescriptor{\n")
	defer func() {
		fmt.Fprintf(logFp, "\t},\n")
	}()

	webservice := TourService{}
	webservice.HTTPHandler().ServeHTTP(recorder, req)

	// dump readable response
	var responseBody bytes.Buffer
	json.Indent(&responseBody, recorder.Body.Bytes(), "", "\t")

	fmt.Fprintf(logFp, "\tStatus:%d,\n", recorder.Code)

	headersToBeSorted = []string{}
	for key, values := range recorder.Header() {
		for _, value := range values {
			headersToBeSorted = append(headersToBeSorted, fmt.Sprintf("%s:%s", key, value))
		}
	}
	sort.Strings(headersToBeSorted)

	fmt.Fprintf(logFp, "\tHeaders:[]string{\n")
	for _, h := range headersToBeSorted {
		fmt.Fprintf(logFp, "\"%s\",\n", h)
	}
	fmt.Fprintf(logFp, "\t},\n")
	fmt.Fprintf(logFp, "\tBody:\n`%s`,\n", responseBody.String())

	return recorder.Code, nil, nil

}
