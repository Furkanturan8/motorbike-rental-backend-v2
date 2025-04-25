package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

var defaultLogger *Logger

func Init(logDir string) error {
	logger, err := NewLogger(logDir)
	if err != nil {
		return err
	}
	defaultLogger = logger
	return nil
}

func NewLogger(logDir string) (*Logger, error) {
	// Log dizinini oluştur
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("log dizini oluşturulamadı: %v", err)
	}

	currentTime := time.Now().Format("2006-01-02")

	// Info log dosyası
	infoLogFile, err := os.OpenFile(
		filepath.Join(logDir, fmt.Sprintf("info_%s.log", currentTime)),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return nil, fmt.Errorf("info log dosyası oluşturulamadı: %v", err)
	}

	// Error log dosyası
	errorLogFile, err := os.OpenFile(
		filepath.Join(logDir, fmt.Sprintf("error_%s.log", currentTime)),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return nil, fmt.Errorf("error log dosyası oluşturulamadı: %v", err)
	}

	return &Logger{
		infoLogger:  log.New(infoLogFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(errorLogFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}, nil
}

func (l *Logger) Info(format string, v ...interface{}) {
	if l.infoLogger != nil {
		l.infoLogger.Printf(format, v...)
	}
	log.Printf(format, v...) // Konsola da yazdır
}

func (l *Logger) Error(format string, v ...interface{}) {
	if l.errorLogger != nil {
		l.errorLogger.Printf(format, v...)
	}
	log.Printf("ERROR: "+format, v...) // Konsola da yazdır
}

// Global fonksiyonlar
func Info(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Info(format, v...)
	}
}

func Error(format string, v ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Error(format, v...)
	}
}
