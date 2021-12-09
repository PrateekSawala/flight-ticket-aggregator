package logging

import (
	"github.com/afiskon/promtail-client/promtail"
	"github.com/sirupsen/logrus"
)

// Logger enforces specific log message formats
type Logger struct {
	Data           map[string]interface{}
	LogrusEntry    *logrus.Entry
	promtailClient promtail.Client
}

// Log is a helper class that enrichens the structured logging
func Log(input interface{}) *Logger {
	data := map[string]interface{}{
		"data": input,
	}
	logrusEntry := logrus.WithFields(logrus.Fields{
		"data": input,
	})
	logger := &Logger{
		Data:        data,
		LogrusEntry: logrusEntry,
	}
	if input == nil {
		return logger
	}

	/*
		// Promtail client
		promtailLabels := fmt.Sprintf("{env=\"%s\", service=\"%s\", tx=\"%d\"}", os.Getenv("FTA_ENVIRONMENT"), os.Getenv("FTA_SERVICE_NAME"), time.Now().Unix())
		// Promtail client config
		conf := promtail.ClientConfig{
			PushURL:            os.Getenv("FTA_LOKI") + "/api/prom/push",
			Labels:             promtailLabels,
			BatchWait:          5 * time.Second,
			BatchEntriesNumber: 20000,
			SendLevel:          promtail.DEBUG,
			PrintLevel:         promtail.ERROR,
		}
		promtailClient, _ := promtail.NewClientProto(conf)
		logger.promtailClient = promtailClient
	*/
	return logger
}
