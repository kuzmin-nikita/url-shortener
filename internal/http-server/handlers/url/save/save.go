package save

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/config"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLSaver
type URLSaver interface {
	SaveURL(urlToSave, alias string) error
}

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

func New(cfg *config.Config, logger *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		logger := logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			logger.Error("Failed to decode request body", sl.Err(err))

			render.JSON(rw, r, resp.Error("Failed to decode request"))

			return
		}

		logger.Info("Request body decoded", slog.Any("request", req))

		if err = validator.New().Struct(req); err != nil {
			validationErros := err.(validator.ValidationErrors)

			logger.Error("Invalid request", sl.Err(err))

			render.JSON(rw, r, resp.ValidationError(validationErros))

			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(cfg.AliasLength)
		}

		// TODO: add check for repetative aliases

		err = urlSaver.SaveURL(req.URL, alias)
		if err != nil {
			if errors.Is(err, storage.ErrorAliasExists) {
				logger.Info("Alias already exists", slog.String("alias", alias))
				render.JSON(rw, r, resp.Error("Alias already exists"))
				return
			}

			logger.Error("Failed to save URL", sl.Err(err))
			render.JSON(rw, r, resp.Error("Failed to save URL"))
			return
		}

		logger.Info("url added")

		render.JSON(rw, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}
