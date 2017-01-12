// Generated automatically by golangAnnotations: do not edit manually

package rest

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/MarcGrol/golangAnnotations/generator/rest/errorh"
	"github.com/gorilla/mux"
)

// HTTPHandler registers endpoint in new router
func (ts *TourService) HTTPHandler() http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	return ts.HTTPHandlerWithRouter(router)
}

// HTTPHandlerWithRouter registers endpoint in existing router
func (ts *TourService) HTTPHandlerWithRouter(router *mux.Router) *mux.Router {
	subRouter := router.PathPrefix("/api/tour").Subrouter()

	subRouter.HandleFunc("/{year}", getTourOnUid(ts)).Methods("GET")

	subRouter.HandleFunc("/{year}/etappe", createEtappe(ts)).Methods("POST")

	subRouter.HandleFunc("/{year}/etappe/{etappeUid}", addEtappeResults(ts)).Methods("PUT")

	subRouter.HandleFunc("/{year}/cyclist", createCyclist(ts)).Methods("POST")

	subRouter.HandleFunc("/{year}/cyclist/{cyclistUid}", markCyclistAbondoned(ts)).Methods("DELETE")

	return router
}

// getTourOnUid does the http handling for business logic method service.getTourOnUid
func getTourOnUid(service *TourService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		pathParams := mux.Vars(r)
		log.Printf("pathParams:%+v", pathParams)

		// extract url-params
		validationErrors := []errorh.FieldError{}

		year := 0

		yearString, exists := pathParams["year"]
		if !exists {

			validationErrors = append(validationErrors, errorh.FieldError{
				SubCode: 1000,
				Field:   "year",
				Msg:     "Missing value for mandatory parameter %s",
				Args:    []string{"year"},
			})

		} else {
			year, err = strconv.Atoi(yearString)
			if err != nil {
				validationErrors = append(validationErrors, errorh.FieldError{
					SubCode: 1001,
					Field:   "year",
					Msg:     "Invalid value for mandatory parameter %s",
					Args:    []string{"year"},
				})
			}
		}

		if len(validationErrors) > 0 {
			errorh.HandleHttpError(errorh.NewInvalidInputErrorSpecific(0, validationErrors), w)
			return
		}

		// call business logic

		result, err := service.getTourOnUid(year)

		if err != nil {
			errorh.HandleHttpError(err, w)
			return
		}

		// write OK response body

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			log.Printf("Error encoding response payload %+v", err)
		}

	}
}

// createEtappe does the http handling for business logic method service.createEtappe
func createEtappe(service *TourService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		pathParams := mux.Vars(r)
		log.Printf("pathParams:%+v", pathParams)

		// extract url-params
		validationErrors := []errorh.FieldError{}

		year := 0

		yearString, exists := pathParams["year"]
		if !exists {

			validationErrors = append(validationErrors, errorh.FieldError{
				SubCode: 1000,
				Field:   "year",
				Msg:     "Missing value for mandatory parameter %s",
				Args:    []string{"year"},
			})

		} else {
			year, err = strconv.Atoi(yearString)
			if err != nil {
				validationErrors = append(validationErrors, errorh.FieldError{
					SubCode: 1001,
					Field:   "year",
					Msg:     "Invalid value for mandatory parameter %s",
					Args:    []string{"year"},
				})
			}
		}

		if len(validationErrors) > 0 {
			errorh.HandleHttpError(errorh.NewInvalidInputErrorSpecific(0, validationErrors), w)
			return
		}

		// read and parse request body
		var etappe Etappe
		err = json.NewDecoder(r.Body).Decode(&etappe)
		if err != nil {
			errorh.HandleHttpError(errorh.NewInvalidInputErrorf(1, "Error parsing request body: %s", err), w)
			return
		}

		// call business logic

		result, err := service.createEtappe(year, etappe)

		if err != nil {
			errorh.HandleHttpError(err, w)
			return
		}

		// write OK response body

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			log.Printf("Error encoding response payload %+v", err)
		}

	}
}

// addEtappeResults does the http handling for business logic method service.addEtappeResults
func addEtappeResults(service *TourService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		pathParams := mux.Vars(r)
		log.Printf("pathParams:%+v", pathParams)

		// extract url-params
		validationErrors := []errorh.FieldError{}

		year := 0

		yearString, exists := pathParams["year"]
		if !exists {

			validationErrors = append(validationErrors, errorh.FieldError{
				SubCode: 1000,
				Field:   "year",
				Msg:     "Missing value for mandatory parameter %s",
				Args:    []string{"year"},
			})

		} else {
			year, err = strconv.Atoi(yearString)
			if err != nil {
				validationErrors = append(validationErrors, errorh.FieldError{
					SubCode: 1001,
					Field:   "year",
					Msg:     "Invalid value for mandatory parameter %s",
					Args:    []string{"year"},
				})
			}
		}

		etappeUid, exists := pathParams["etappeUid"]
		if !exists {

			validationErrors = append(validationErrors, errorh.FieldError{
				SubCode: 1000,
				Field:   "etappeUid",
				Msg:     "Missing value for mandatory parameter %s",
				Args:    []string{"etappeUid"},
			})

		}

		if len(validationErrors) > 0 {
			errorh.HandleHttpError(errorh.NewInvalidInputErrorSpecific(0, validationErrors), w)
			return
		}

		// read and parse request body
		var results EtappeResult
		err = json.NewDecoder(r.Body).Decode(&results)
		if err != nil {
			errorh.HandleHttpError(errorh.NewInvalidInputErrorf(1, "Error parsing request body: %s", err), w)
			return
		}

		// call business logic

		err = service.addEtappeResults(year, etappeUid, results)

		if err != nil {
			errorh.HandleHttpError(err, w)
			return
		}

		// write OK response body

		w.WriteHeader(http.StatusNoContent)

	}
}

// createCyclist does the http handling for business logic method service.createCyclist
func createCyclist(service *TourService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		pathParams := mux.Vars(r)
		log.Printf("pathParams:%+v", pathParams)

		// extract url-params
		validationErrors := []errorh.FieldError{}

		year := 0

		yearString, exists := pathParams["year"]
		if !exists {

			validationErrors = append(validationErrors, errorh.FieldError{
				SubCode: 1000,
				Field:   "year",
				Msg:     "Missing value for mandatory parameter %s",
				Args:    []string{"year"},
			})

		} else {
			year, err = strconv.Atoi(yearString)
			if err != nil {
				validationErrors = append(validationErrors, errorh.FieldError{
					SubCode: 1001,
					Field:   "year",
					Msg:     "Invalid value for mandatory parameter %s",
					Args:    []string{"year"},
				})
			}
		}

		if len(validationErrors) > 0 {
			errorh.HandleHttpError(errorh.NewInvalidInputErrorSpecific(0, validationErrors), w)
			return
		}

		// read and parse request body
		var cyclist Cyclist
		err = json.NewDecoder(r.Body).Decode(&cyclist)
		if err != nil {
			errorh.HandleHttpError(errorh.NewInvalidInputErrorf(1, "Error parsing request body: %s", err), w)
			return
		}

		// call business logic

		result, err := service.createCyclist(year, cyclist)

		if err != nil {
			errorh.HandleHttpError(err, w)
			return
		}

		// write OK response body

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(result)
		if err != nil {
			log.Printf("Error encoding response payload %+v", err)
		}

	}
}

// markCyclistAbondoned does the http handling for business logic method service.markCyclistAbondoned
func markCyclistAbondoned(service *TourService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		pathParams := mux.Vars(r)
		log.Printf("pathParams:%+v", pathParams)

		// extract url-params
		validationErrors := []errorh.FieldError{}

		year := 0

		yearString, exists := pathParams["year"]
		if !exists {

			validationErrors = append(validationErrors, errorh.FieldError{
				SubCode: 1000,
				Field:   "year",
				Msg:     "Missing value for mandatory parameter %s",
				Args:    []string{"year"},
			})

		} else {
			year, err = strconv.Atoi(yearString)
			if err != nil {
				validationErrors = append(validationErrors, errorh.FieldError{
					SubCode: 1001,
					Field:   "year",
					Msg:     "Invalid value for mandatory parameter %s",
					Args:    []string{"year"},
				})
			}
		}

		cyclistUid, exists := pathParams["cyclistUid"]
		if !exists {

			validationErrors = append(validationErrors, errorh.FieldError{
				SubCode: 1000,
				Field:   "cyclistUid",
				Msg:     "Missing value for mandatory parameter %s",
				Args:    []string{"cyclistUid"},
			})

		}

		if len(validationErrors) > 0 {
			errorh.HandleHttpError(errorh.NewInvalidInputErrorSpecific(0, validationErrors), w)
			return
		}

		// call business logic

		err = service.markCyclistAbondoned(year, cyclistUid)

		if err != nil {
			errorh.HandleHttpError(err, w)
			return
		}

		// write OK response body

		w.WriteHeader(http.StatusNoContent)

	}
}
