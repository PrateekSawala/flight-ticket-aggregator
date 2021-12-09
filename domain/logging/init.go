package logging

import (
	"github.com/sirupsen/logrus"
)

// Initialize Logging sets the logging environment
func InitializeLogging() {
	logrus.SetLevel(logrus.TraceLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableQuote: true,
	})
}
