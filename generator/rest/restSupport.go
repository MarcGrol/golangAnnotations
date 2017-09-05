package rest

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"

	"github.com/MarcGrol/golangAnnotations/generator/rest/errorh"
)

type restErrorHandler interface {
	HandleRestError(c context.Context, credentials Credentials, error errorh.Error, r *http.Request)
}

var RestErrorHandler restErrorHandler

func HandleHttpError(c context.Context, credentials Credentials, err error, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(errorh.GetHttpCode(err))

	errorResp := errorh.MapToError(err)

	if RestErrorHandler != nil {
		RestErrorHandler.HandleRestError(c, credentials, errorResp, r)
	}

	// write response
	json.NewEncoder(w).Encode(errorResp)
}

type Credentials struct {
	Language      string
	RequestUID    string
	SessionUID    string
	EndUserAccess string
	EndUserRole   string
	EndUserUID    string
	ApiKey        string
}

func ExtractCredentials(language string, r *http.Request) Credentials {
	return Credentials{
		Language:      language,
		RequestUID:    r.Header.Get("X-request-uid"),
		SessionUID:    r.Header.Get("X-session-uid"),
		EndUserAccess: r.Header.Get("X-enduser-access"),
		EndUserRole:   r.Header.Get("X-enduser-role"),
		EndUserUID:    r.Header.Get("X-enduser-uid"),
	}
}
