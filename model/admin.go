package model

import (
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminSuccessLogin struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
	Token *Token `json:"token"`
}

type AdminRegistration struct {
	Name     string `validate:"required" json:"name"`
	Email    string `validate:"required,email" json:"email"`
	Password string `validate:"required" json:"password"`
	Role     string `validate:"required" json:"role"`
}

type AdminLogin struct {
	Email    string `validate:"required,email" json:"email"`
	Password string `validate:"required" json:"password"`
}

type AdminBearer struct {
	jwt.StandardClaims
	ID              primitive.ObjectID `json:"id"`
	Name            string             `json:"name"`
	Email           string             `json:"email"`
	EmailIsVerified bool               `json:"emailIsVerified"`
	ExpiresAt       int64              `json:"exp,omitempty"`
}

type Token struct {
	Value     *string `json:"value,omitempty"`
	ExpiresIn int64   `json:"expiresIn"`
}
