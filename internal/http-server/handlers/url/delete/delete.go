package delete

import (
	"log/slog"
	"net/http"

	resp "url-shortener/internal/lib/api/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLDeletter
type URLDeletter interface {
	DeleteURL(alias string) error
}

func New(logger *slog.Logger, urlDeletter URLDeletter) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		logger := logger.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")

		if alias == "" {
			logger.Error("alias is empty")
			render.JSON(rw, r, resp.Error("alias is empty"))
			return
		}

		err := urlDeletter.DeleteURL(alias)
		if err != nil {
			logger.Error("failed to delete URL")
			render.JSON(rw, r, resp.Error("failed to delete URL"))
			return
		}

		logger.Info("url deleted", slog.String("alias", alias))
	}
}
