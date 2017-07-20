package rest

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/net/context"
)

type Credentials struct {
	Language      string
	RequestUID    string
	SessionUID    string
	EndUserAccess string
	EndUserRole   string
	EndUserUID    string
	ApiKey        string
}

type restSupport interface {
	Logger

	CreateContext(r *http.Request) context.Context

	HandleHttpError(c context.Context, credentials Credentials, err error, w http.ResponseWriter, r *http.Request)
}

var Support restSupport

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
