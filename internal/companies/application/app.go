package application

import (
	"encoding/json"
	"net/http"

	domain "github.com/erotokritosVall/xmapp/internal/companies/domain"
	"github.com/erotokritosVall/xmapp/pkg/api"
	"github.com/erotokritosVall/xmapp/pkg/errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

type companiesApp struct {
	srv       domain.CompanyService
	validator *validator.Validate
}

func NewApp(srv domain.CompanyService, validator *validator.Validate) api.App {
	return &companiesApp{
		srv:       srv,
		validator: validator,
	}
}

func (app *companiesApp) RegisterPublicEndpoints(router chi.Router) {
	router.Get("/v1/companies/{id}", app.read)
}

func (app *companiesApp) RegisterProtectedEndpoints(router chi.Router) {
	router.Post("/v1/companies", app.insert)
	router.Patch("/v1/companies/{id}", app.update)
	router.Delete("/v1/companies/{id}", app.delete)
}

func (app *companiesApp) read(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")

	company, err := app.srv.Read(request.Context(), id)
	if err != nil {
		if err == errors.ErrNotFound {
			log.Warn().Err(err).Msgf("could not find company %s", id)
			api.Response(http.StatusNotFound).Send(writer, request)
			return
		}

		log.Err(err).Msgf("could not find company %s", id)
		api.Response(http.StatusInternalServerError).
			WithError(errors.ErrGeneric).
			Send(writer, request)
		return
	}

	resp := companyToReadDto(company)

	api.Response(http.StatusOK).
		WithBody(resp, nil).
		Send(writer, request)
}

func (app *companiesApp) insert(writer http.ResponseWriter, request *http.Request) {
	req := &insertCompanyRequest{}
	if err := json.NewDecoder(request.Body).Decode(req); err != nil {
		log.Warn().Err(err).Msg("failed to decode insertCompanyRequest")
		api.Response(http.StatusBadRequest).Send(writer, request)
		return
	}

	if err := app.validator.Struct(req); err != nil {
		log.Warn().Err(err).Msg("failed to validate insertCompanyRequest")

		api.Response(http.StatusBadRequest).
			WithError(err).
			Send(writer, request)
		return
	}

	company, err := req.toDomain()
	if err != nil {
		log.Warn().Err(err).Msg("failed to create company")

		api.Response(http.StatusBadRequest).
			WithError(err).
			Send(writer, request)
		return
	}

	id, err := app.srv.Insert(request.Context(), company)
	if err != nil {
		log.Err(err).Msg("failed to insert company")

		api.Response(http.StatusInternalServerError).
			WithError(errors.ErrGeneric).
			Send(writer, request)
		return
	}

	api.Response(http.StatusCreated).
		WithBody(id, nil).
		Send(writer, request)
}

func (app *companiesApp) update(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")

	req := &updateCompanyRequest{}
	if err := json.NewDecoder(request.Body).Decode(req); err != nil {
		log.Warn().Err(err).Msg("failed to decode updateCompanyRequest")
		api.Response(http.StatusBadRequest).Send(writer, request)
		return
	}

	opts := req.toDomain()

	err := app.srv.Update(request.Context(), id, opts)
	if err != nil {
		if err == errors.ErrNotFound {
			log.Warn().Err(err).Msgf("could not find company %s", id)
			api.Response(http.StatusNotFound).Send(writer, request)
			return
		}

		log.Err(err).Msgf("failed to update company %s", id)

		api.Response(http.StatusInternalServerError).
			WithError(errors.ErrGeneric).
			Send(writer, request)
		return
	}

	api.Response(http.StatusOK).
		Send(writer, request)
}

func (app *companiesApp) delete(writer http.ResponseWriter, request *http.Request) {
	id := chi.URLParam(request, "id")

	if err := app.srv.Delete(request.Context(), id); err != nil {
		log.Err(err).Msgf("failed to delete company %s", id)

		api.Response(http.StatusInternalServerError).
			WithError(errors.ErrGeneric).
			Send(writer, request)
		return
	}

	api.Response(http.StatusOK).Send(writer, request)
}
