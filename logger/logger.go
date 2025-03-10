package logger

import (
	"errors"
	"io"
	"os"
	"runtime/debug"
	"sync"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

// Logger is the main logging struct
var (
	Logger   *logrus.Logger
	logMutex sync.Mutex
)

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

// Logging functions
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

// Send logs to external logging services
func AddLogSink(writer io.Writer) error {
	if writer == nil {
		return errors.New("invalid log writer")
	}
	Logger.SetOutput(io.MultiWriter(Logger.Out, writer))
	return nil
}
