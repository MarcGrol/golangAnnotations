// Generated automatically: do not edit manually

package testData

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/MarcGrol/microgen/lib/myerrors"
	"github.com/gorilla/mux"
)

func (ts *MyService) HttpHandler() http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/person", doit(ts)).Methods("GET")

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
