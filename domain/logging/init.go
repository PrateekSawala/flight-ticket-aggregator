package logging

import (
	"github.com/sirupsen/logrus"
)

// Initialize Logging sets the logging environment
func InitializeLogging() {
	logrus.SetLevel(logrus.TraceLevel)
}
