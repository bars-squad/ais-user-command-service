package jwt

import (
	"context"
	"crypto/rsa"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

// Errors.
var (
	ErrInvalidToken      error = fmt.Errorf("invalid token")
	ErrExpiredOrNotReady error = fmt.Errorf("token is either expired or not ready to use")
)

// JSONWebToken is a concrete struct of json web token.
type JSONWebToken struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

// NewJSONWebToken is a constructor.
func NewJSONWebToken(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) *JSONWebToken {
	return &JSONWebToken{privateKey, publicKey}
}

// Sign will generate new jwt token.
func (a *JSONWebToken) Sign(ctx context.Context, claims jwt.Claims) (tokenString string, err error) {
	// span, _ := apm.StartSpan(ctx, "JSONWebToken: Sign", "token.jwt")
	// defer span.End()

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(a.privateKey)
}

// Parse will parse the token string to bearer claims.
func (a *JSONWebToken) Parse(ctx context.Context, tokenString string, claims jwt.Claims) (err error) {
	// span, _ := apm.StartSpan(ctx, "JSONWebToken: Parse", "token.jwt")
	// defer span.End()

	token, err := jwt.ParseWithClaims(tokenString, claims, a.keyFunc)
	if err = a.checkError(err); err != nil {
		return
	}

	if !token.Valid {
		return ErrInvalidToken
	}

	return
}

func (a *JSONWebToken) keyFunc(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, ErrInvalidToken
	}
	return a.publicKey, nil
}

func (a *JSONWebToken) checkError(err error) error {
	if err == nil {
		return err
	}

	ve, ok := err.(*jwt.ValidationError)
	if !ok {
		return ErrInvalidToken
	}
	if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
		return ErrExpiredOrNotReady
	}

	return ErrInvalidToken
}
