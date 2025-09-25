package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	baseLogger   *logrus.Logger
	defaultEntry *logrus.Entry
)

// Init configures a logrus logger.
func Init(serviceName string) {
	if baseLogger != nil {
		return
	}

	l := logrus.New()
	l.SetOutput(os.Stdout)
	l.SetLevel(logrus.InfoLevel)
	l.SetReportCaller(false)
	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
	})

	baseLogger = l
	defaultEntry = l.WithField("service", serviceName)
}

// L returns the default logger entry enriched with base fields.
func L() *logrus.Entry {
	if defaultEntry == nil {
		Init("app")
	}
	return defaultEntry
}

// With returns a new entry with additional fields.
func With(fields logrus.Fields) *logrus.Entry {
	return L().WithFields(fields)
}
