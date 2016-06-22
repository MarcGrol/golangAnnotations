package testData

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MarcGrol/microgen/lib/myerrors"
)

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
