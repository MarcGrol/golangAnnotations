package rest

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"

	"google.golang.org/appengine/user"

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

// @JsonStruct()
type Credentials struct {
	Language      string    `json:"language,omitempty"`
	RequestUID    string    `json:"requestUid,omitempty"`
	SessionUID    string    `json:"sessionUid,omitempty"`
	EndUserAccess string    `json:"endUserAccess,omitempty"`
	EndUserRole   string    `json:"endUserRole,omitempty"`
	EndUserUID    string    `json:"endUserUid,omitempty"`
	ApiKey        string    `json:"apiKey,omitempty"`
	AuthUser      *AuthUser `json:"authUser,omitempty"`
}

// provided by App Engine's user authentication service.
// @JsonStruct()
type AuthUser struct {
	EmailAddress string `json:"emailAddress,omitempty"`
	AuthDomain   string `json:"authDomain,omitempty"`
	IsAdmin      bool   `json:"isAdmin,omitempty"`
	ID           string `json:"id,omitempty"`
}

func ExtractAllCredentials(c context.Context, r *http.Request, language string) Credentials {
	return Credentials{
		Language:      language,
		RequestUID:    r.Header.Get("X-request-uid"),
		SessionUID:    r.Header.Get("X-session-uid"),
		EndUserAccess: r.Header.Get("X-enduser-access"),
		EndUserRole:   r.Header.Get("X-enduser-role"),
		EndUserUID:    r.Header.Get("X-enduser-uid"),
		AuthUser:      GetAuthUser(c),
	}
}

func ExtractAdminCredentials(c context.Context, r *http.Request, language string) Credentials {
	return Credentials{
		AuthUser: GetAuthUser(c),
	}
}

func GetAuthUser(c context.Context) *AuthUser {
	u := user.Current(c)
	if u == nil {
		return nil
	}
	return &AuthUser{
		EmailAddress: u.Email,
		AuthDomain:   u.AuthDomain,
		IsAdmin:      u.Admin,
		ID:           u.ID,
	}
}
