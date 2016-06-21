// Generated automatically: do not edit manually

package web

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/MarcGrol/microgen/lib/myerrors"
	"github.com/gorilla/mux"
)

func (ts *TourService) HandleHttp() http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/tour/{year}", getTourOnUid(ts)).Methods("GET")

	router.HandleFunc("/tour/{year}/etappe", createEtappe(ts)).Methods("POST")

	router.HandleFunc("/tour/{year}/etappe/{etappeUid}", addEtappeResults(ts)).Methods("POST")

	router.HandleFunc("/tour/{year}/cyclist", createCyclist(ts)).Methods("POST")

	router.HandleFunc("/tour/{year}/cyclist/{cyclistUid}", markCyclistAbondoned(ts)).Methods("DELETE")

	return router
}

func getTourOnUid(service *TourService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		pathParams := mux.Vars(r)
		log.Printf("pathParams:%+v", pathParams)

		// extract url-params

		yearString, exists := pathParams["year"]
		if !exists {
			handleError(myerrors.NewInvalidInputError(fmt.Errorf("Missing path param 'year'")), w)
			return
		}
		year, err := strconv.Atoi(yearString)
		if err != nil {
			handleError(myerrors.NewInvalidInputError(fmt.Errorf("Invalid path param 'year'")), w)
			return
		}

		// call business logic

		result, err := service.getTourOnUid(year)

		if err != nil {
			handleError(err, w)
			return
		}

		// write response body

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			log.Printf("Error encoding response payload %+v", err)
		}

	}
}

func createEtappe(service *TourService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		pathParams := mux.Vars(r)
		log.Printf("pathParams:%+v", pathParams)

		// extract url-params

		yearString, exists := pathParams["year"]
		if !exists {
			handleError(myerrors.NewInvalidInputError(fmt.Errorf("Missing path param 'year'")), w)
			return
		}
		year, err := strconv.Atoi(yearString)
		if err != nil {
			handleError(myerrors.NewInvalidInputError(fmt.Errorf("Invalid path param 'year'")), w)
			return
		}

		// read abd parse request body
		var etappe Etappe
		err = json.NewDecoder(r.Body).Decode(&etappe)
		if err != nil {
			handleError(myerrors.NewInvalidInputError(fmt.Errorf("Error decoding request payload:%s", err)), w)
			return
		}

		// call business logic

		result, err := service.createEtappe(year, etappe)

		if err != nil {
			handleError(err, w)
			return
		}

		// write response body

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			log.Printf("Error encoding response payload %+v", err)
		}

	}
}

func addEtappeResults(service *TourService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		pathParams := mux.Vars(r)
		log.Printf("pathParams:%+v", pathParams)

		// extract url-params

		yearString, exists := pathParams["year"]
		if !exists {
			handleError(myerrors.NewInvalidInputError(fmt.Errorf("Missing path param 'year'")), w)
			return
		}
		year, err := strconv.Atoi(yearString)
		if err != nil {
			handleError(myerrors.NewInvalidInputError(fmt.Errorf("Invalid path param 'year'")), w)
			return
		}

		etappeUid, exists := pathParams["etappeUid"]
		if !exists {
			handleError(myerrors.NewInvalidInputError(fmt.Errorf("Missing path param 'etappeUid'")), w)
			return
		}

		// read abd parse request body
		var results EtappeResult
		err = json.NewDecoder(r.Body).Decode(&results)
		if err != nil {
			handleError(myerrors.NewInvalidInputError(fmt.Errorf("Error decoding request payload:%s", err)), w)
			return
		}

		// call business logic

		err = service.addEtappeResults(year, etappeUid, results)

		if err != nil {
			handleError(err, w)
			return
		}

		// write response body

		w.WriteHeader(http.StatusNoContent)

	}
}

func createCyclist(service *TourService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		pathParams := mux.Vars(r)
		log.Printf("pathParams:%+v", pathParams)

		// extract url-params

		yearString, exists := pathParams["year"]
		if !exists {
			handleError(myerrors.NewInvalidInputError(fmt.Errorf("Missing path param 'year'")), w)
			return
		}
		year, err := strconv.Atoi(yearString)
		if err != nil {
			handleError(myerrors.NewInvalidInputError(fmt.Errorf("Invalid path param 'year'")), w)
			return
		}

		// read abd parse request body
		var cyclist Cyclist
		err = json.NewDecoder(r.Body).Decode(&cyclist)
		if err != nil {
			handleError(myerrors.NewInvalidInputError(fmt.Errorf("Error decoding request payload:%s", err)), w)
			return
		}

		// call business logic

		result, err := service.createCyclist(year, cyclist)

		if err != nil {
			handleError(err, w)
			return
		}

		// write response body

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			log.Printf("Error encoding response payload %+v", err)
		}

	}
}

func markCyclistAbondoned(service *TourService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		pathParams := mux.Vars(r)
		log.Printf("pathParams:%+v", pathParams)

		// extract url-params

		yearString, exists := pathParams["year"]
		if !exists {
			handleError(myerrors.NewInvalidInputError(fmt.Errorf("Missing path param 'year'")), w)
			return
		}
		year, err := strconv.Atoi(yearString)
		if err != nil {
			handleError(myerrors.NewInvalidInputError(fmt.Errorf("Invalid path param 'year'")), w)
			return
		}

		cyclistUid, exists := pathParams["cyclistUid"]
		if !exists {
			handleError(myerrors.NewInvalidInputError(fmt.Errorf("Missing path param 'cyclistUid'")), w)
			return
		}

		// call business logic

		err = service.markCyclistAbondoned(year, cyclistUid)

		if err != nil {
			handleError(err, w)
			return
		}

		// write response body

		w.WriteHeader(http.StatusNoContent)

	}
}
