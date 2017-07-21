package rest

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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

	errorResp := errorh.Error{
		ErrorMessage: err.Error(),
		ErrorCode:    errorh.GetErrorCode(err),
		FieldErrors:  errorh.GetFieldErrors(err),
	}

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
	username, password, err := decodeBasicAuthHeader(r)
	if err == nil {
		return Credentials{
			Language:      language,
			RequestUID:    r.Header.Get("X-request-uid"),
			SessionUID:    "",
			EndUserAccess: "",
			EndUserRole:   "supplier",
			EndUserUID:    username,
			ApiKey:        password,
		}
	}
	return Credentials{
		Language:      language,
		RequestUID:    r.Header.Get("X-request-uid"),
		SessionUID:    r.Header.Get("X-session-uid"),
		EndUserAccess: r.Header.Get("X-enduser-access"),
		EndUserRole:   r.Header.Get("X-enduser-role"),
		EndUserUID:    r.Header.Get("X-enduser-uid"),
	}
}

func decodeBasicAuthHeader(r *http.Request) (string, string, error) {
	authHeader := r.Header["Authorization"]
	if len(authHeader) == 0 {
		return "", "", fmt.Errorf("Missing header")
	}

	auth := strings.SplitN(authHeader[0], " ", 2)
	if len(auth) != 2 || auth[0] != "Basic" {
		return "", "", fmt.Errorf("Invalid header")
	}

	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", 2)
	if len(pair) != 2 {
		return "", "", fmt.Errorf("Invalid/missing header-values")
	}

	return pair[0], pair[1], nil
}
