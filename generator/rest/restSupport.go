package rest

import "net/http"

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

type Credentials struct {
	Language 	  string
	RequestUID    string
	SessionUID    string
	EndUserAccess string
	EndUserRole   string
	EndUserUID    string
}
