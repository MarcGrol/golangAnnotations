// Generated automatically: do not edit manually

package testData

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/MarcGrol/microgen/lib/myerrors"
	"github.com/gorilla/mux"
)

func (ts *MyService) HttpHandler() http.Handler {
	servicePrefix := "/api"
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc(servicePrefix+"/person", doit(ts)).Methods("GET")

	return router
}

func doit(service *MyService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		pathParams := mux.Vars(r)
		log.Printf("pathParams:%+v", pathParams)

		// extract url-params

		uidString, exists := pathParams["uid"]
		if !exists {
			handleError(myerrors.NewInvalidInputError(fmt.Errorf("Missing path param 'uid'")), w)
			return
		}
		uid, err := strconv.Atoi(uidString)
		if err != nil {
			handleError(myerrors.NewInvalidInputError(fmt.Errorf("Invalid path param 'uid'")), w)
			return
		}

		subuid, exists := pathParams["subuid"]
		if !exists {
			handleError(myerrors.NewInvalidInputError(fmt.Errorf("Missing path param 'subuid'")), w)
			return
		}

		// call business logic

		err = service.doit(uid, subuid)

		if err != nil {
			handleError(err, w)
			return
		}

		// write response body

		w.WriteHeader(http.StatusNoContent)

	}
}

func handleError(err error, w http.ResponseWriter) {
	errorBody := struct {
		ErrorMessage string
	}{
		err.Error(),
	}
	blob, err := json.Marshal(errorBody)
	if err != nil {
		log.Printf("Error marshalling error response payload %+v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(determineHttpCode(err))
	w.Header().Set("Content-Type", "application/json")
	w.Write(blob)
}

func determineHttpCode(err error) int {
	if myerrors.IsNotFoundError(err) {
		return http.StatusNotFound
	} else if myerrors.IsInternalError(err) {
		return http.StatusInternalServerError
	} else if myerrors.IsInvalidInputError(err) {
		return http.StatusBadRequest
	} else if myerrors.IsNotAuthorizedError(err) {
		return http.StatusForbidden
	} else {
		return http.StatusInternalServerError
	}
}
