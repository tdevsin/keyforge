package logger

import (
	"sync"

	"go.uber.org/zap"
)

var (
	log  *zap.Logger
	once sync.Once
)

// getLogger returns the global logger instance. If the logger has not been initialized yet, it will create a new logger instance.
func GetLogger() *zap.Logger {
	once.Do(func() {
		if log != nil {
			// If the logger is already set, use it
			return
		}
		var err error
		baseLogger, err := zap.NewProduction() // Replace with zap.NewDevelopment() for dev env
		if err != nil {
			log = zap.NewNop() // Fallback to no-op logger
			return
		}
		log = baseLogger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))
	})
	return log
}

// Info logs a message at the info level.
func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// Error logs a message at the error level.
func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

// Debug logs a message at the debug level.
func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

// Warn logs a message at the warn level.
func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

// Sync flushes any buffered log entries. This should be called before exiting the program.
func Sync() {
	logger := GetLogger()
	if logger != nil {
		_ = logger.Sync()
	}
}
