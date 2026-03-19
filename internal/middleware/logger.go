package middleware

import (
	"atur-dana/internal/metrics"
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		metrics.RequestsInFlight.Inc()
		defer metrics.RequestsInFlight.Dec()

		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		metrics.RecordRequest(r.Method, r.URL.Path, rw.status, duration.Seconds())

		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.status,
			"latency_ms", duration.Milliseconds(),
			"request_id", GetRequestID(r),
		)
	})
}
