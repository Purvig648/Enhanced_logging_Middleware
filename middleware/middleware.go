package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/Purvig648/Enhanced_logging_Middleware/logger"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// ContextKey is a custom type to avoid key collisions in context
type ContextKey string

const TraceIDKey ContextKey = "trace_id"

// GenerateTraceID creates a new trace ID
func GenerateTraceID() string {
	return uuid.New().String()
}

// MiddlewareLogger logs HTTP requests and assigns a trace ID, also logs response status
func MiddlewareLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		traceID := r.Header.Get("X-Trace-ID")
		if traceID == "" {
			traceID = GenerateTraceID()
		}

		ctx := context.WithValue(r.Context(), TraceIDKey, traceID)
		r = r.WithContext(ctx)

		logger.Logger.WithFields(logrus.Fields{
			"trace_id":  traceID,
			"method":    r.Method,
			"path":      r.URL.Path,
			"remote_ip": r.RemoteAddr,
		}).Info("Incoming request")

		statusRecorder := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(statusRecorder, r)

		logger.Logger.WithFields(logrus.Fields{
			"trace_id": traceID,
			"duration": time.Since(start).String(),
			"status":   statusRecorder.statusCode,
		}).Info("Request completed")
	})
}

// responseWriter is a wrapper to capture HTTP status codes
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
