package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/difaal21/go-template/entity"
	"github.com/difaal21/go-template/jwt"
	"github.com/difaal21/go-template/model"
	"github.com/difaal21/go-template/responses"
)

var (
	header = "Authorization"
)

const (
	invalidTokenMessage = "Invalid token"
)

// Session is concrete struct of jwt authorization.
type Session struct {
	// sess         session.Session
	jsonWebToken *jwt.JSONWebToken
}

// NewSessionMiddleware is a constructor.
// func NewSessionMiddleware(sess session.Session, jsonWebToken *jwt.JSONWebToken) RouteMiddleware {
// 	return &Session{sess, jsonWebToken}
// }

// NewSessionMiddleware is a constructor.
func NewSessionMiddleware(jsonWebToken *jwt.JSONWebToken) RouteMiddleware {
	return &Session{jsonWebToken}
}

// Verify will verify the http incomming request to ensure it comes within the authorized token.
func (a Session) Verify(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// defer span.End()

		authHeader := r.Header.Get(header)
		if authHeader == "" {
			a.respondUnauthorized(w, jwt.ErrInvalidToken.Error())
			// span.End()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 {
			a.respondUnauthorized(w, jwt.ErrInvalidToken.Error())
			// span.End()
			return
		}

		tokenString := bearerToken[1]

		var userBearerClaims model.AdminBearer

		if err := a.jsonWebToken.Parse(ctx, tokenString, &userBearerClaims); err != nil {
			a.respondUnauthorized(w, err.Error())
			// span.End()
			return
		}

		ctx = context.WithValue(ctx, &entity.AdminContextKey{}, &userBearerClaims)
		r = r.WithContext(ctx)

		/* userBuff, err := a.sess.Get(ctx, userBearerClaims.NationalityId)
		if err != nil {
			a.respondUnauthorized(w, err.Error())
			// span.End()
			return
		} */

		// user := entity.User{}
		// json.Unmarshal(userBuff, &user)
		// ctx = context.WithValue(ctx, entity.UserContextKey{}, user)

		// r = r.WithContext(ctx)

		// span.End()
		next(w, r)
	})
}

func (a Session) respondUnauthorized(w http.ResponseWriter, message string) {
	responses.REST(w, httpResponse.Unathorized("").NewResponses(nil, invalidTokenMessage))
}
