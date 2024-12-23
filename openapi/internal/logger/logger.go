package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	errorLogger *log.Logger
	infoLogger  *log.Logger
)

// InitLogger 初始化日志记录器
func InitLogger() error {
	// 创建日志目录
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %v", err)
	}

	// 创建当天的日志文件
	currentTime := time.Now()
	errorLogPath := filepath.Join(logDir, fmt.Sprintf("error_%s.log", currentTime.Format("2006-01-02")))
	infoLogPath := filepath.Join(logDir, fmt.Sprintf("info_%s.log", currentTime.Format("2006-01-02")))

	// 打开错误日志文件
	errorFile, err := os.OpenFile(errorLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open error log file: %v", err)
	}

	// 打开信息日志文件
	infoFile, err := os.OpenFile(infoLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to open info log file: %v", err)
	}

	// 创建多重写入器，同时写入文件和标准输出
	errorWriter := io.MultiWriter(errorFile, os.Stdout)
	infoWriter := io.MultiWriter(infoFile, os.Stdout)

	// 初始化日志记录器
	errorLogger = log.New(errorWriter, "\033[31m[ERROR]\033[0m ", log.Ldate|log.Ltime|log.Lshortfile)
	infoLogger = log.New(infoWriter, "\033[32m[INFO]\033[0m ", log.Ldate|log.Ltime)

	return nil
}

// Error 记录错误日志
func Error(format string, v ...interface{}) {
	if errorLogger != nil {
		errorLogger.Printf(format, v...)
	}
}

// Info 记录信息日志
func Info(format string, v ...interface{}) {
	if infoLogger != nil {
		infoLogger.Printf(format, v...)
	}
}
