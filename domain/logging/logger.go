package logging

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// Log is a helper class that enrichens the structured logging
func Log(input interface{}) *logrus.Entry {
	log := logrus.WithFields(logrus.Fields{
		"method": "unkown",
	})
	// Check input
	if input == nil {
		return log
	}

	// Check input type
	switch input.(type) {
	case string:
		return log.WithField("method", input)
	default:
		// Add raw input
		rawInput := fmt.Sprintf("%+v", input)
		if rawInput != "<nil>" {
			log = log.WithField("input", rawInput)
		}
		return log.WithField("method", input)
	}
}
