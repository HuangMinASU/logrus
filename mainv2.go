package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// CustomHook defines a custom logrus Hook for the desensitization process.
type CustomHook struct{}

// Levels returns the logging levels applicable for this hook.
func (hook *CustomHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire processes and logs each entry's message when a log event occurs.
func (hook *CustomHook) Fire(entry *logrus.Entry) error {
	output, err := FormatJSON(entry.Message)
	if err != nil {
		output = entry.Message // Fallback to the original message
	}

	entry.Message = output
	fmt.Println(output) // Print to console

	return nil
}

// FormatJSON formats the message and desensitizes specific fields.
func FormatJSON(message string) (string, error) {
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(message), &jsonData); err != nil {
		return fmt.Sprintf("[String] %s", message), nil
	}

	// Desensitization of specific fields
	desensitizeData(jsonData, "keyToModify")

	compactJSON, err := json.Marshal(jsonData)
	if err != nil {
		return message, err
	}
	return string(compactJSON), nil
}

func desensitizeData(data map[string]interface{}, key string) {
	if value, exists := data[key].(string); exists && len(value) > 6 {
		data[key] = value[:len(value)-6] + "******"
	}

	for _, v := range data {
		if nestedMap, ok := v.(map[string]interface{}); ok {
			desensitizeData(nestedMap, key)
		}
	}
}

func main() {
	// Initialize loggers
	logA := logrus.New()
	logB := logrus.New()
	logC := logrus.New()

	// Define logKindToFileName map
	logKindToFileName := map[string]*logrus.Logger{
		"./log/access/mpcmediaaccessservice-api-access.log": logA,
		"./log/access/run.log":                              logB,
		"./log/access/DEBUG.log":                            logC,
	}

	// Set output files for each logger using the mapping
	for fileName, logger := range logKindToFileName {
		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			logger.SetOutput(file)
		} else {
			logger.Info("Failed to log to file, using default stdout")
		}

		// Configure formatter
		standardFormatter := &logrus.TextFormatter{
			FullTimestamp: true,
		}
		logger.SetFormatter(standardFormatter)
	}

	// Add a custom hook for desensitization to loggers A and B
	hook := &CustomHook{}
	logA.AddHook(hook)
	logB.AddHook(hook)

	// Example logging
	jsonMessage := `{"nestedField": {"keyToModify": "SensitiveValue"}, "anotherField": "value"}`
	logKindToFileName["./log/access/mpcmediaaccessservice-api-access.log"].Info(jsonMessage)

	jsonMessage2 := `{"nestedField": {"keyToModify": "22543322555"}, "anotherField": "value"}`
	logKindToFileName["./log/access/run.log"].Info(jsonMessage2)

	logKindToFileName["./log/access/DEBUG.log"].Info("This is a debug message.")
}
