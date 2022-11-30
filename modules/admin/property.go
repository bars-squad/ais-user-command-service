package admin

import (
	"github.com/bars-squad/ais-user-command-service/jwt"
	"github.com/bars-squad/ais-user-command-service/session"
	"github.com/sirupsen/logrus"
)

type Property struct {
	ServiceName string
	Logger      *logrus.Logger
	Repository  Repository
	// SubsidyProductRepository subsidyproduct.Repository
	JSONWebToken *jwt.JSONWebToken
	Session      session.Session
	// Publisher                pubsub.Publisher
}
