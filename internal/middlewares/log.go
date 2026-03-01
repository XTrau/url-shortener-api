package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (srw *responseWriter) WriteHeader(statusCode int) {
	srw.statusCode = statusCode
	srw.ResponseWriter.WriteHeader(statusCode)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				slog.Error(fmt.Sprintf("Error recovered: %v", r))
			}
		}()

		start := time.Now()
		rw := &responseWriter{w, 200}

		next.ServeHTTP(rw, r)

		elapsed := time.Since(start)
		slog.Info(
			"request",
			slog.String("Method", r.Method),
			slog.String("Path", r.URL.Path),
			slog.Int("Status", rw.statusCode),
			slog.Duration("Duration", elapsed),
		)
	})
}
