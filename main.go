package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

func main() {
	// 创建一个新的logger实例
	logger := logrus.New()

	// 设置日志格式，可以使用JSONFormatter或TextFormatter
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

	// 设置日志级别
	logger.SetLevel(logrus.InfoLevel)

	// 打印不同级别的日志
	logger.Info("This is an info message")
	logger.Warn("This is a warning message")
	logger.Error("This is an error message")

	// 使用带有字段的日志
	logger.WithFields(logrus.Fields{
		"event": "event1",
		"topic": "topic1",
	}).Info("Logging with Fields")

	// 封装一个函数来处理高级日志记录
	logWithFields(logger, "event2", "topic2", "This is another log with fields")
}

func logWithFields(logger *logrus.Logger, event, topic, message string) {
	logger.WithFields(logrus.Fields{
		"event": event,
		"topic": topic,
	}).Info(message)
}
