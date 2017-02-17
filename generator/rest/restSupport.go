package rest

import "net/http"

func ExtractAuthContext(language string, r *http.Request) map[string]string {
	return map[string]string{
		"language":      language,
		"requestUid":    r.Header.Get("X-request-uid"),
		"sessionUid":    r.Header.Get("X-session-uid"),
		"enduserAccess": r.Header.Get("X-enduser-access"),
		"enduserRole":   r.Header.Get("X-enduser-role"),
		"enduserUid":    r.Header.Get("X-enduser-uid"),
	}
}
