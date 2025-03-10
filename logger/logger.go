package logger

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

// Logger is the main logging struct
var (
	Logger   *logrus.Logger
	logMutex sync.Mutex
)

// ContextKey is a custom type to avoid key collisions in context
type ContextKey string

const TraceIDKey ContextKey = "trace_id"

// InitLogger initializes the Logrus logger with configurable settings and log rotation
func InitLogger(format string, level string, logFile string, maxSize int, maxBackups int, maxAge int) {
	Logger = logrus.New()

	// Set log output (stdout or file with rotation)
	if logFile != "" {
		Logger.Out = &lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    maxSize,    // megabytes
			MaxBackups: maxBackups, // number of backups
			MaxAge:     maxAge,     // days
			Compress:   true,
		}
	} else {
		Logger.Out = os.Stdout
	}

	// Set log format (JSON or text)
	if format == "json" {
		Logger.SetFormatter(&logrus.JSONFormatter{})
	} else {
		Logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	// Set log level
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.InfoLevel
	}
	Logger.SetLevel(lvl)
}

// Log functions with structured error logging
func Info(message string, fields map[string]interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()
	Logger.WithFields(fields).Info(message)
}

func Error(message string, fields map[string]interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()
	fields["stack_trace"] = string(debug.Stack())
	Logger.WithFields(fields).Error(message)
}

func Debug(message string, fields map[string]interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()
	Logger.WithFields(fields).Debug(message)
}

func Warn(message string, fields map[string]interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()
	Logger.WithFields(fields).Warn(message)
}

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

		Logger.WithFields(logrus.Fields{
			"trace_id":  traceID,
			"method":    r.Method,
			"path":      r.URL.Path,
			"remote_ip": r.RemoteAddr,
		}).Info("Incoming request")

		statusRecorder := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(statusRecorder, r)

		Logger.WithFields(logrus.Fields{
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

// Send logs to external logging services
func AddLogSink(writer io.Writer) error {
	if writer == nil {
		return errors.New("invalid log writer")
	}
	Logger.SetOutput(io.MultiWriter(Logger.Out, writer))
	return nil
}

// AsyncLog handles log writing asynchronously to improve performance
func AsyncLog(level logrus.Level, message string, fields map[string]interface{}) {
	go func() {
		logMutex.Lock()
		defer logMutex.Unlock()
		Logger.WithFields(fields).Log(level, message)
	}()
}
