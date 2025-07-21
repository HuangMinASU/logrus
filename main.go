package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// 是否美化输出的全局开关
var prettyPrint = false

// FormatJSON 根据全局开关决定是否美化输出 JSON
func FormatJSON(message string) (string, error) {
	var jsonData map[string]interface{}
	if err := json.Unmarshal([]byte(message), &jsonData); err != nil {
		return fmt.Sprintf("[String] %s", message), nil
	}

	if prettyPrint {
		prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")
		if err != nil {
			return message, err
		}
		return string(prettyJSON), nil
	} else {
		compactJSON, err := json.Marshal(jsonData)
		if err != nil {
			return message, err
		}
		return string(compactJSON), nil
	}
}

// CustomHook 定义一个自定义的 Hook 结构体
type CustomHook struct{}

// Levels 返回这个 Hook 所应用的日志级别
func (hook *CustomHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire 在日志发生时被调用，对日志的 Message 进行处理和打印
func (hook *CustomHook) Fire(entry *logrus.Entry) error {
	output, err := FormatJSON(entry.Message)
	if err != nil {
		output = entry.Message // 回退到原始信息
	}

	// 打印输出
	fmt.Println(output)

	return nil
}

func main() {
	// 创建一个新的 logger 实例
	logger := logrus.New()

	// 创建并添加 CustomHook 实例
	hook := &CustomHook{}
	logger.AddHook(hook)

	// 设置日志格式，可以使用 JSONFormatter 或 TextFormatter
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true, // 包含完整的时间戳
	})

	// 设置日志输出位置，可以是标准输出或文件
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.SetOutput(file)
	} else {
		logger.Info("Failed to log to file, using default stdout")
	}

	// 打印包含 JSON 数据的日志
	jsonMessage := `{"nestedField": {"keyToModify": "OldValue"}, "anotherField": "value"}`
	logger.Info(jsonMessage)

	// 打印普通字符串的日志
	plainMessage := "This is a plain text message."
	logger.Info(plainMessage)

	// 使用带有字段的日志
	logger.WithFields(logrus.Fields{
		"event": "event1",
		"topic": "topic1",
	}).Info("Logging with Fields")
}
