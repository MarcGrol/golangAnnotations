// Generated automatically: do not edit manually

package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
)

func getTourOnUidTestHelper(url string) (int, *Tour, error) {

	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {

		return 0, nil, err

	}

	req.Header.Set("Accept", "application/json")

	webservice := TourService{}
	webservice.HttpHandler().ServeHTTP(recorder, req)

	var resp Tour
	dec := json.NewDecoder(recorder.Body)
	err = dec.Decode(&resp)
	if err != nil {
		return recorder.Code, nil, err
	}
	return recorder.Code, &resp, nil

}

func createEtappeTestHelper(url string, input Etappe) (int, *Etappe, error) {

	recorder := httptest.NewRecorder()

	requestBody, _ := json.Marshal(input)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(requestBody)))

	if err != nil {

		return 0, nil, err

	}

	req.Header.Set("Accept", "application/json")

	webservice := TourService{}
	webservice.HttpHandler().ServeHTTP(recorder, req)

	var resp Etappe
	dec := json.NewDecoder(recorder.Body)
	err = dec.Decode(&resp)
	if err != nil {
		return recorder.Code, nil, err
	}
	return recorder.Code, &resp, nil

}

func addEtappeResultsTestHelper(url string, input EtappeResult) (int, error) {

	recorder := httptest.NewRecorder()

	requestBody, _ := json.Marshal(input)
	req, err := http.NewRequest("PUT", url, strings.NewReader(string(requestBody)))

	if err != nil {

		return 0, err

	}

	webservice := TourService{}
	webservice.HttpHandler().ServeHTTP(recorder, req)

	return recorder.Code, nil

}

func createCyclistTestHelper(url string, input Cyclist) (int, *Cyclist, error) {

	recorder := httptest.NewRecorder()

	requestBody, _ := json.Marshal(input)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(requestBody)))

	if err != nil {

		return 0, nil, err

	}

	req.Header.Set("Accept", "application/json")

	webservice := TourService{}
	webservice.HttpHandler().ServeHTTP(recorder, req)

	var resp Cyclist
	dec := json.NewDecoder(recorder.Body)
	err = dec.Decode(&resp)
	if err != nil {
		return recorder.Code, nil, err
	}
	return recorder.Code, &resp, nil

}

func markCyclistAbondonedTestHelper(url string) (int, error) {

	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("DELETE", url, nil)

	if err != nil {

		return 0, err

	}

	webservice := TourService{}
	webservice.HttpHandler().ServeHTTP(recorder, req)

	return recorder.Code, nil

}
