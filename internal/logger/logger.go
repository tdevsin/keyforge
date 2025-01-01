package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logging interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Sync()
}

type Logger struct {
	Logger *zap.Logger
}

// getLogger returns the global logger instance. If the logger has not been initialized yet, it will create a new logger instance.
func GetLogger(isProd bool, nodeId string) *Logger {
	var logger Logger
	var err error
	var baseLogger *zap.Logger
	if isProd {
		config := zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000 MST")
		baseLogger, err = config.Build()
	} else {
		baseLogger, err = zap.NewDevelopment()
	}
	if err != nil {
		logger.Logger = zap.NewNop() // Fallback to no-op logger
		return &logger
	}
	logger.Logger = baseLogger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).With(zap.String("nodeId", nodeId))
	return &logger
}

// Info logs a message at the info level.
func (l *Logger) Info(msg string, fields ...zap.Field) {
	l.Logger.Info(msg, fields...)
}

// Error logs a message at the error level.
func (l *Logger) Error(msg string, fields ...zap.Field) {
	l.Logger.Error(msg, fields...)
}

// Debug logs a message at the debug level.
func (l *Logger) Debug(msg string, fields ...zap.Field) {
	l.Logger.Debug(msg, fields...)
}

// Warn logs a message at the warn level.
func (l *Logger) Warn(msg string, fields ...zap.Field) {
	l.Logger.Warn(msg, fields...)
}

// Sync flushes any buffered log entries. This should be called before exiting the program.
func (l *Logger) Sync() {
	l.Logger.Sync()
}

// Infof formats and logs a message at the info level, satisfying Pebble's Logger interface.
func (l *Logger) Infof(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.Logger.Info(message)
}

// Fatalf formats and logs a message at the fatal level and then terminates the application, satisfying Pebble's Logger interface.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.Logger.Fatal(message)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.Logger.Error(message)
}
