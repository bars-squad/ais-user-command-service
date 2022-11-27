package entity

import (
	"context"
	"fmt"
	"time"

	"github.com/difaal21/go-template/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminContextKey struct{}

func GetAdministratorFromContext(ctx context.Context) (admin *model.AdminBearer, err error) {
	administratorSessionValue := ctx.Value(&AdminContextKey{})
	admin, ok := administratorSessionValue.(*model.AdminBearer)
	if !ok {
		err = fmt.Errorf("invalid administrator context key")
		return
	}

	return
}

type Admin struct {
	ID              primitive.ObjectID `json:"id" bson:"_id"`
	Name            string             `json:"name" bson:"name"`
	Email           string             `json:"email" bson:"email"`
	EmailIsVerified bool               `json:"emailIsVerified" bson:"emailIsVerified"`
	Password        any                `json:"password" bson:"password"`
	Role            string             `json:"role" bson:"role"`
	CreatedBy       CreatedBy          `json:"createdBy,omitempty" bson:"createdBy"`
	CreatedAt       *time.Time         `json:"createdAt" bson:"createdAt"`
	UpdatedAt       *time.Time         `json:"updatedAt" bson:"updatedAt"`
}

type CreatedBy struct {
	UserID primitive.ObjectID `json:"id" bson:"userId"`
	Name   string             `json:"name" bson:"name"`
	Email  string             `json:"email" bson:"email"`
}
