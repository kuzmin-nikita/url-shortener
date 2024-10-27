package mwLogger

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log := log.With(slog.String("component", "middleware/logger"))

		log.Info("Logger middleware enabled")

		fn := func(responseWriter http.ResponseWriter, request *http.Request) {
			entry := log.With(
				slog.String("method", request.Method),
				slog.String("path", request.URL.Path),
				slog.String("remote_addr", request.RemoteAddr),
				slog.String("request_id", middleware.GetReqID(request.Context())),
			)
			wrapResponseWriter := middleware.NewWrapResponseWriter(responseWriter, request.ProtoMajor)

			startTime := time.Now()

			defer func() {
				entry.Info("Request completed",
					slog.Int("status", wrapResponseWriter.Status()),
					slog.Int("bytes", wrapResponseWriter.BytesWritten()),
					slog.String("duration", time.Since(startTime).String()),
				)
			}()

			next.ServeHTTP(wrapResponseWriter, request)
		}

		return http.HandlerFunc(fn)
	}
}
