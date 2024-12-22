package logger

import "github.com/sirupsen/logrus"

type LogrusEntryImpl struct {
	entry *logrus.Entry
}

func (l *LogrusEntryImpl) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

func (l *LogrusEntryImpl) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

func (l *LogrusEntryImpl) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

func (l *LogrusEntryImpl) Info(args ...interface{}) {
	l.entry.Info(args...)
}

func (l *LogrusEntryImpl) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

func (l *LogrusEntryImpl) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

func (l *LogrusEntryImpl) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

func (l *LogrusEntryImpl) Error(args ...interface{}) {
	l.entry.Error(args...)
}

func (l *LogrusEntryImpl) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

func (l *LogrusEntryImpl) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

func (l *LogrusEntryImpl) WithField(key string, value interface{}) (entry Logger) {
	entry = &LogrusEntryImpl{l.entry.WithField(key, value)}
	return
}

func (l *LogrusEntryImpl) WithFields(args map[string]interface{}) (entry Logger) {
	entry = &LogrusEntryImpl{l.entry.WithFields(args)}
	return
}
