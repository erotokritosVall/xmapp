package application

import (
	"encoding/json"
	"net/http"

	domain "github.com/erotokritosVall/xmapp/internal/users/domain"
	"github.com/erotokritosVall/xmapp/pkg/api"
	"github.com/erotokritosVall/xmapp/pkg/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

type userApp struct {
	srv       domain.UserService
	validator *validator.Validate
}

func NewApp(srv domain.UserService) api.App {
	return &userApp{
		srv: srv,
	}
}

func (app *userApp) RegisterProtectedEndpoints(router chi.Router) {
	router.Post("/v1/logout", app.logout)
	router.Post("/v1/users", app.createUser)
}

func (app *userApp) RegisterPublicEndpoints(router chi.Router) {
	router.Post("/v1/login", app.login)
}

func (app *userApp) login(writer http.ResponseWriter, request *http.Request) {
	req := &loginRequest{}
	if err := json.NewDecoder(request.Body).Decode(req); err != nil {
		log.Warn().Msgf("failed to decode loginRequest: %+v", err)

		api.Response(http.StatusBadRequest).
			WithError(err).
			Send(writer, request)
		return
	}

	if err := app.validator.Struct(req); err != nil {
		log.Warn().Msgf("failed to validate loginRequest: %+v", err)

		api.Response(http.StatusBadRequest).
			WithError(err).
			Send(writer, request)
		return
	}

	token, err := app.srv.Login(request.Context(), req.Email, req.Password)
	if err != nil {
		log.Warn().Err(err).Msgf("user %s failed to login", req.Email)

		api.Response(http.StatusBadRequest).
			WithError(errors.ErrGeneric).
			Send(writer, request)
		return
	}

	api.Response(http.StatusOK).
		WithBody(token, nil).
		Send(writer, request)
}

func (app *userApp) logout(writer http.ResponseWriter, request *http.Request) {
	if err := app.srv.Logout(request.Context()); err != nil {
		log.Warn().Msgf("failed to logout: %+v", err)

		api.Response(http.StatusInternalServerError).
			WithError(err).
			Send(writer, request)
		return
	}

	api.Response(http.StatusOK).Send(writer, request)
}

func (app *userApp) createUser(writer http.ResponseWriter, request *http.Request) {
	req := &createUserRequest{}
	if err := json.NewDecoder(request.Body).Decode(req); err != nil {
		log.Warn().Msgf("failed to decode createUserRequest: %+v", err)

		api.Response(http.StatusBadRequest).
			WithError(err).
			Send(writer, request)
		return
	}

	if err := app.validator.Struct(req); err != nil {
		log.Warn().Msgf("failed to validate createUserRequest: %+v", err)

		api.Response(http.StatusBadRequest).
			WithError(err).
			Send(writer, request)
		return
	}

	u := &domain.User{
		Email:    req.Email,
		Password: req.Password,
	}

	id, err := app.srv.Create(request.Context(), u)
	if err != nil {
		log.Error().Err(err).Msgf("failed to create user %s", req.Email)

		api.Response(http.StatusInternalServerError).
			WithError(errors.ErrGeneric).
			Send(writer, request)
		return
	}

	api.Response(http.StatusCreated).
		WithBody(id, nil).
		Send(writer, request)
}
