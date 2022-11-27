package admin

import (
	"encoding/json"
	"net/http"

	"github.com/difaal21/go-template/helpers/validation"
	"github.com/difaal21/go-template/middleware"
	"github.com/difaal21/go-template/model"
	"github.com/difaal21/go-template/responses"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	unprocessableEntityMessage = "Invalid payload format"
	internalServerErrorMessage = "Internal server error"
	badRequestMessage          = "Please check your payload"
)

var (
	httpResponse = responses.HttpResponseStatusCodesImpl{}
)

type HTTPHandler struct {
	Logger   *logrus.Logger
	Validate *validator.Validate
	Usecase  Usecase
}

func NewHTTPHandler(logger *logrus.Logger, validate *validator.Validate, router *mux.Router, basicAuth middleware.RouteMiddleware, usecase Usecase, sess middleware.RouteMiddleware) {
	handler := &HTTPHandler{
		Logger:   logger,
		Validate: validate,
		Usecase:  usecase,
	}

	router.HandleFunc("/v1/admin/login", basicAuth.Verify(handler.Login)).Methods(http.MethodPost)
	router.HandleFunc("/v1/admin/registration", sess.Verify(handler.Register)).Methods(http.MethodPost)
	// router.HandleFunc("/mpv-general-registration/v1/users/registration/{nationalityId}", basicAuth.Verify(handler.GetUser)).Methods(http.MethodGet)
	// router.HandleFunc("/mpv-general-registration/v1/users/registration/{nationalityId}/subsidy-product/{subsidyProduct}", basicAuth.Verify(handler.GetUserByNationalityIDAndSubsidyProduct)).Methods(http.MethodGet)
}

func (handler *HTTPHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload *model.AdminLogin

	ctx := r.Context()
	defer func() {
		r := recover()
		if r != nil {
			handler.Logger.Error(r)
			responses.REST(w, httpResponse.InternalServerError("").NewResponses(nil, internalServerErrorMessage))
		}
	}()

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.REST(w, httpResponse.UnprocessableEntity("").NewResponses(nil, err.Error()))
		return
	}

	if err := validation.RequestBody(handler.Validate, payload); err != nil {
		responses.REST(w, httpResponse.BadRequest("").NewResponses(err, badRequestMessage))
		return
	}

	resp := handler.Usecase.Login(ctx, payload)
	responses.REST(w, resp)
}

func (handler *HTTPHandler) Register(w http.ResponseWriter, r *http.Request) {
	var payload *model.AdminRegistration

	ctx := r.Context()
	defer func() {
		r := recover()
		if r != nil {
			handler.Logger.Error(r)
			responses.REST(w, httpResponse.InternalServerError("").NewResponses(nil, internalServerErrorMessage))
		}
	}()

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		responses.REST(w, httpResponse.UnprocessableEntity("").NewResponses(nil, err.Error()))
		return
	}

	if err := validation.RequestBody(handler.Validate, payload); err != nil {
		responses.REST(w, httpResponse.BadRequest("").NewResponses(err, badRequestMessage))
		return
	}

	var role, err = RoleValidation(payload.Role)
	if err != nil {
		payload.Role = role.String()

		handler.Logger.Error(err)
		responses.REST(w, httpResponse.Forbidden("").NewResponses(nil, err.Error()))
		return
	}

	resp := handler.Usecase.Register(ctx, payload)
	responses.REST(w, resp)
}
