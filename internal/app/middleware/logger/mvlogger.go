package mwlogger

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
)

func NewMwLogger(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		log.Info("middleware logger")

		logFn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			newW := middleware.NewWrapResponseWriter(w, r.ProtoMajor) //меняем стандартный Response методом из chi
			next.ServeHTTP(newW, r)

			duration := time.Since(start)

			log.Info("request completed",
				slog.String("uri", r.RequestURI),
				slog.String("method", r.Method),
				slog.Duration("duration", duration),
				slog.Int("status", newW.Status()),
				slog.Int("size", newW.BytesWritten()),
			)

		}
		return http.HandlerFunc(logFn)
	}
}
