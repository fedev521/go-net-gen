package log

import (
	"log"
	"os"

	"github.com/sirupsen/logrus"
	logrusadapter "logur.dev/adapter/logrus"
	"logur.dev/logur"
)

// NewErrorStandardLogger returns a new standard logger logging on error level.
func NewErrorStandardLogger(logger logur.Logger) *log.Logger {
	return logur.NewErrorStandardLogger(logger, "", 0)
}

// SetStandardLogger sets the global logger's output to a custom logger
// instance.
func SetStandardLogger(logger logur.Logger) {
	log.SetOutput(logur.NewLevelWriter(logger, logur.Info))
}

// NewSimpleLogger returns a new logger that requires no configuration. It
// writes in JSON format to stdout at Info level.
func NewSimpleLogger() logur.LoggerFacade {
	logger := logrus.New()

	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)

	return logrusadapter.New(logger)
}

// NewLogger creates a new logger with the given configuration. It is resilient
// to malformed or invalid log configuration.
func NewLogger(config Config) logur.LoggerFacade {
	logger := logrus.New()

	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{})

	level, err := logrus.ParseLevel(config.Level)
	if err == nil {
		logger.SetLevel(level)
	}

	return logrusadapter.New(logger)
}
