package log

import (
	"log"
	"os"
)

// Loggerはログ出力を管理します
type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

// NewLoggerは新しいLoggerインスタンスを作成します
func NewLogger() *Logger {
	return &Logger{
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Infoは情報ログを出力します
func (l *Logger) Info(v ...interface{}) {
	l.infoLogger.Println(v...)
}

// Errorはエラーログを出力します
func (l *Logger) Error(v ...interface{}) {
	l.errorLogger.Println(v...)
}
