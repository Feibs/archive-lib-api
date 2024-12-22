package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type LogrusLoggerImpl struct {
	log *logrus.Logger
}

func NewLogrusLogger() *LogrusLoggerImpl {
	log := logrus.New()

	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		PrettyPrint:     false,
	})

	log.SetLevel(logrus.InfoLevel)
	log.SetOutput(os.Stdout)

	return &LogrusLoggerImpl{
		log: log,
	}
}

func (l *LogrusLoggerImpl) Debug(args ...any) {
	l.log.Debug(args...)
}

func (l *LogrusLoggerImpl) Debugf(format string, args ...any) {
	l.log.Debugf(format, args...)
}

func (l *LogrusLoggerImpl) Info(args ...any) {
	l.log.Info(args...)
}

func (l *LogrusLoggerImpl) Infof(format string, args ...any) {
	l.log.Infof(format, args...)
}

func (l *LogrusLoggerImpl) Warn(args ...any) {
	l.log.Warn(args...)
}

func (l *LogrusLoggerImpl) Warnf(format string, args ...any) {
	l.log.Warnf(format, args...)
}

func (l *LogrusLoggerImpl) Error(args ...any) {
	l.log.Error(args...)
}

func (l *LogrusLoggerImpl) Errorf(format string, args ...any) {
	l.log.Errorf(format, args...)
}

func (l *LogrusLoggerImpl) Fatal(args ...any) {
	l.log.Fatal(args...)
}

func (l *LogrusLoggerImpl) Fatalf(format string, args ...any) {
	l.log.Fatalf(format, args...)
}

func (l *LogrusLoggerImpl) WithField(key string, value any) Logger {
	return &LogrusEntryImpl{
		entry: l.log.WithField(key, value),
	}
}

func (l *LogrusLoggerImpl) WithFields(fields map[string]any) Logger {
	return &LogrusEntryImpl{
		entry: l.log.WithFields(fields),
	}
}
