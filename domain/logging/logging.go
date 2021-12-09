package logging

import (
	// "fmt"
	"github.com/sirupsen/logrus"
)

// Trace is a wrapper for the logrus Trace method
func (l *Logger) Trace(args ...interface{}) {
	l.LogrusEntry.Trace(args...)
	// l.promtailClient.Debugf(fmt.Sprintf("%+v", args...))
}

// Debug is a wrapper for the logrus Trace method
func (l *Logger) Debug(args ...interface{}) {
	l.LogrusEntry.Debug(args...)
	// l.promtailClient.Debugf(fmt.Sprintf("%+v", args...))
}

// Info is a wrapper for the logrus Trace method
func (l *Logger) Info(args ...interface{}) {
	l.LogrusEntry.Info(args...)
	// l.promtailClient.Infof(fmt.Sprintf("%+v", args...))
}

// Warn is a wrapper for the logrus Trace method
func (l *Logger) Warn(args ...interface{}) {
	l.LogrusEntry.Warn(args...)
	// l.promtailClient.Warnf(fmt.Sprintf("%+v", args...))
}

// Error is a wrapper for the logrus Trace method
func (l *Logger) Error(args ...interface{}) {
	l.LogrusEntry.Error(args...)
	// l.promtailClient.Errorf(fmt.Sprintf("%+v", args...))
}

// Tracef is a wrapper for the logrus Tracef method
func (l *Logger) Tracef(format string, args ...interface{}) {
	l.LogrusEntry.Tracef(format, args...)
	// l.promtailClient.Debugf(format, args...)
}

// Debugf is a wrapper for the logrus Debugf method
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.LogrusEntry.Debugf(format, args...)
	// l.promtailClient.Debugf(format, args...)
}

// Infof is a wrapper for the logrus Infof method
func (l *Logger) Infof(format string, args ...interface{}) {
	l.LogrusEntry.Infof(format, args...)
	// l.promtailClient.Infof(format, args...)
}

// Warnf is a wrapper for the logrus Warnf method
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.LogrusEntry.Warnf(format, args...)
	// l.promtailClient.Warnf(format, args...)
}

// Errorf is a wrapper for the logrus Errorf method
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.LogrusEntry.Errorf(format, args...)
	// l.promtailClient.Errorf(format, args...)
}

// WithFields add the provided fields to the logger
func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.LogrusEntry.WithFields(fields)
}
