package middleware

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Auth provides the middleware to authenticate against
// the registered API key.
type Auth struct {
	apiKey          string
	apiKeyParamName string
}

// NewAuth instantiates our middleware object.
func NewAuth(apiKey string, apiKeyParamName string) *Auth {
	return &Auth{apiKey: apiKey, apiKeyParamName: apiKeyParamName}
}

// ServeHTTP deals with authenticating against the registered API key.
// The API key is expected to be in the query string where one is provided.
func (a *Auth) ServeHTTP(w http.ResponseWriter, r *http.Request, ps httprouter.Params, next http.HandlerFunc) {
	apiKey := r.URL.Query().Get(a.apiKeyParamName)
	if apiKey == a.apiKey {
		next(w, r)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("{\"message\":\"You are not authorised to access this webhook endpoint\"}"))
	}
}
