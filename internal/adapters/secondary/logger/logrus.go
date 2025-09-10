package logger

import (
	"logistics-api/internal/pkg/logger"

	"github.com/sirupsen/logrus"
)

type LogrusAdapter struct {
	logger *logrus.Logger
}

func NewLogrusAdapter(level string, format string) *LogrusAdapter {
	log := logrus.New()

	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	if format == "json" {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	return &LogrusAdapter{logger: log}
}

func (l *LogrusAdapter) Debug(msg string, fields ...logger.Field) {
	l.logger.WithFields(l.convertFields(fields)).Debug(msg)
}

func (l *LogrusAdapter) Info(msg string, fields ...logger.Field) {
	l.logger.WithFields(l.convertFields(fields)).Info(msg)
}

func (l *LogrusAdapter) Warn(msg string, fields ...logger.Field) {
	l.logger.WithFields(l.convertFields(fields)).Warn(msg)
}

func (l *LogrusAdapter) Error(msg string, fields ...logger.Field) {
	l.logger.WithFields(l.convertFields(fields)).Error(msg)
}

func (l *LogrusAdapter) Fatal(msg string, fields ...logger.Field) {
	l.logger.WithFields(l.convertFields(fields)).Fatal(msg)
}

func (l *LogrusAdapter) With(fields ...logger.Field) logger.Logger {
	return &LogrusAdapter{
		logger: l.logger.WithFields(l.convertFields(fields)).Logger,
	}
}

func (l *LogrusAdapter) convertFields(fields []logger.Field) logrus.Fields {
	logrusFields := make(logrus.Fields)
	for _, field := range fields {
		logrusFields[field.Key] = field.Value
	}
	return logrusFields
}
