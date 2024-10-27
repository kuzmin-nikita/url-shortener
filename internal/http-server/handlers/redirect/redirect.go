package redirect

import (
	"errors"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v2@v2.40.2 --name=URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(logger *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.New"

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

		url, err := urlGetter.GetURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrorAliasNotFound) {
				logger.Error("alias not found", slog.String("alias", alias))
				render.JSON(rw, r, resp.Error("alias not found"))
				return
			}

			logger.Error("failed to get URL", sl.Err(err))
			render.JSON(rw, r, resp.Error("failed to get URL"))
			return
		}

		logger.Info("url got", slog.String("url", url))

		http.Redirect(rw, r, url, http.StatusFound)
	}
}
