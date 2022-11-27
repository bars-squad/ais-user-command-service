package admin

import (
	"context"
	"time"

	"github.com/difaal21/go-template/entity"
	"github.com/difaal21/go-template/exception"
	"github.com/difaal21/go-template/helpers/cryptography"
	"github.com/difaal21/go-template/helpers/date"
	"github.com/difaal21/go-template/jwt"
	"github.com/difaal21/go-template/model"
	"github.com/difaal21/go-template/responses"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	httpResponses = responses.HttpResponseStatusCodesImpl{}
)

const (
	duplicatedEmailStatus = "DUPLICATED_EMAIL"
)

const (
	userAlreadyExistMessage   = "User already exist"
	createAdminSuccessMessage = "Successfully created an account"
	invalidCredentialMessage  = "Invalid admin credential"
	loginSucessMessage        = "Login success"
)

type Usecase interface {
	Register(ctx context.Context, param *model.AdminRegistration) responses.Responses
	// Register(ctx context.Context, param model.UserRegistration) (resp response.Response)
	Login(ctx context.Context, param *model.AdminLogin) responses.Responses
	// GetUserProfile(ctx context.Context) (resp response.Response)
	// SetUserEmailActivation(ctx context.Context, param model.UserEmailActivation) (resp response.Response)
	// ChangePassword(ctx context.Context, param model.UserChangePassword) (resp response.Response)
	// UploadImage(ctx context.Context, folderName string, param model.UploadImageParams) (resp response.Response)
	// UpdateProfile(ctx context.Context, param model.UpdateProfileParams) (resp response.Response)
	// ForgotPassword(ctx context.Context, param model.ForgotPasswordParams) (resp response.Response)
	// ResetPassword(ctx context.Context, param model.ResetPasswordParams) (resp response.Response)
	// UpdateEmail(ctx context.Context, param model.UpdateEmailParams) (resp response.Response)
}

type UsecaseImpl struct {
	ServiceName string
	Logger      *logrus.Logger
	Repository  Repository
	// SubsidyProductRepository subsidyproduct.Repository
	JSONWebToken *jwt.JSONWebToken
	// Session                  session.Session
	// Publisher                pubsub.Publisher
}

func NewUsecase(property *Property) Usecase {
	return &UsecaseImpl{
		ServiceName:  property.ServiceName,
		Logger:       property.Logger,
		Repository:   property.Repository,
		JSONWebToken: property.JSONWebToken,
		// sess:                     property.Session,
		// publisher:                property.Publisher,
	}
}

func (u UsecaseImpl) Register(ctx context.Context, param *model.AdminRegistration) responses.Responses {
	var admin entity.Admin

	createdBy, err := entity.GetAdministratorFromContext(ctx)
	if err != nil {
		u.Logger.Error(err.Error())
		return httpResponses.InternalServerError("").NewResponses(nil, internalServerErrorMessage)
	}

	isAdminExist, err := u.Repository.FindOneByEmail(ctx, param.Email)
	if err != nil {
		if err != exception.ErrNotFound {
			u.Logger.Error(err.Error())
			return httpResponses.InternalServerError("").NewResponses(nil, internalServerErrorMessage)
		}
	}

	if isAdminExist != nil {
		u.Logger.Error(userAlreadyExistMessage)
		return httpResponses.Conflict(duplicatedEmailStatus).NewResponses(nil, userAlreadyExistMessage)
	}

	hashPassword, err := cryptography.Hash([]byte(param.Password))
	if err != nil {
		u.Logger.Error(err.Error())
		return httpResponses.InternalServerError("").NewResponses(nil, internalServerErrorMessage)
	}

	admin.ID = primitive.NewObjectID()
	admin.Name = param.Name
	admin.Email = param.Email
	admin.EmailIsVerified = false
	admin.CreatedAt = date.CurrentUTCTime()
	admin.Password = hashPassword
	admin.Role = param.Role
	admin.CreatedBy.UserID = createdBy.ID
	admin.CreatedBy.Name = createdBy.Name
	admin.CreatedBy.Email = createdBy.Email

	if err = u.Repository.Save(ctx, &admin); err != nil {
		u.Logger.Error(err)
		return httpResponses.InternalServerError("").NewResponses(nil, internalServerErrorMessage)
	}

	admin.Password = nil

	return httpResponses.Ok("").NewResponses(admin, createAdminSuccessMessage)
}

func (u UsecaseImpl) Login(ctx context.Context, param *model.AdminLogin) responses.Responses {
	var token model.Token
	var profile model.AdminSuccessLogin

	var admin, err = u.Repository.FindOneByEmail(ctx, param.Email)
	if err != nil {
		if err != exception.ErrNotFound {
			u.Logger.Error(err.Error())
			return httpResponses.InternalServerError("").NewResponses(nil, internalServerErrorMessage)
		}
	}

	if admin == nil {
		param.Password = ""
		u.Logger.WithFields(logrus.Fields{"body": param}).Error(err.Error())
		return httpResponses.BadRequest("").NewResponses(nil, invalidCredentialMessage)
	}

	passwordMatch := cryptography.Verify(admin.Password.(string), []byte(param.Password))
	if !passwordMatch {
		return httpResponses.BadRequest("").NewResponses(nil, invalidCredentialMessage)
	}

	expiresAt := time.Now().Add(time.Hour * 24 * 7).Unix()
	claims := &model.AdminBearer{}
	claims.ID = admin.ID
	claims.Name = admin.Name
	claims.Email = admin.Email
	claims.ExpiresAt = expiresAt

	tokenString, err := u.JSONWebToken.Sign(ctx, claims)
	if err != nil {
		u.Logger.WithFields(logrus.Fields{"body": param}).Error(err.Error())
		return httpResponses.InternalServerError("").NewResponses(nil, err.Error())
	}

	token.Value = &tokenString
	token.ExpiresIn = 60 * 60 * 24 * 7

	profile.Email = admin.Email
	profile.Name = admin.Name
	profile.Role = admin.Role
	profile.Token = &token

	return httpResponses.Ok("").NewResponses(profile, loginSucessMessage)
}
