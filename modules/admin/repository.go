package admin

import (
	"context"

	"github.com/difaal21/go-template/databases/mongodb"
	"github.com/difaal21/go-template/entity"
	"github.com/difaal21/go-template/exception"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository is a collection of behavior of User repository.
type Repository interface {
	Save(ctx context.Context, admin *entity.Admin) (err error)
	FindOneByEmail(ctx context.Context, email string) (admin *entity.Admin, err error)
}

type RepositoryImpl struct {
	logger *logrus.Logger
	col    mongodb.Collection
}

func NewRepository(logger *logrus.Logger, db mongodb.Database) Repository {
	col := db.Collection("administrator")

	return &RepositoryImpl{logger, col}
}

func (repo *RepositoryImpl) Save(ctx context.Context, admin *entity.Admin) (err error) {
	_, err = repo.col.InsertOne(ctx, &admin)
	if err != nil {
		repo.logger.Error(err)
		err = exception.ErrInternalServer
		return
	}

	return
}

func (repo *RepositoryImpl) FindOneByEmail(ctx context.Context, email string) (admin *entity.Admin, err error) {
	filter := bson.M{
		"email": email,
	}

	if err = repo.col.FindOne(ctx, filter).Decode(&admin); err != nil {
		if err != mongo.ErrNoDocuments {
			repo.logger.Error(err)
			err = exception.ErrInternalServer
			return
		}
		err = exception.ErrNotFound
		return
	}

	return
}
