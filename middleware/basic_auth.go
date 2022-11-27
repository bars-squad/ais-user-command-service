package middleware

import (
	"net/http"

	"github.com/difaal21/go-template/responses"
)

const (
	errorMessage = "Invalid token"
)

var (
	httpResponse = responses.HttpResponseStatusCodesImpl{}
)

// BasicAuth is a concrete struct of basic auth verifier.
type BasicAuth struct {
	username, password string
}

// NewBasicAuth is a constructor.
func NewBasicAuth(username, password string) RouteMiddleware {
	return &BasicAuth{username, password}
}

func (ba *BasicAuth) respondUnauthorized(w http.ResponseWriter) {
	responses.REST(w, httpResponse.Unathorized("").NewResponses(nil, errorMessage))
}

// Verify will verify the request to ensure it comes with an authorized basic auth token.
func (ba *BasicAuth) Verify(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// span, _ := apm.StartSpan(r.Context(), "Basic Auth Middleware: Verify", "middleware.basic_auth")

		username, password, ok := r.BasicAuth()
		if !ok {
			ba.respondUnauthorized(w)
			// span.End()
			return
		}

		if !(username == ba.username && password == ba.password) {
			ba.respondUnauthorized(w)
			// span.End()
			return
		}
		// span.End()
		next(w, r)
	})
}
