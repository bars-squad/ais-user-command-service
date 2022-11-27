package admin

import (
	"github.com/difaal21/go-template/jwt"
	"github.com/sirupsen/logrus"
)

type Property struct {
	ServiceName string
	Logger      *logrus.Logger
	Repository  Repository
	// SubsidyProductRepository subsidyproduct.Repository
	JSONWebToken *jwt.JSONWebToken
	// Session                  session.Session
	// Publisher                pubsub.Publisher
}
