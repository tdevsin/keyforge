package logger

import (
	"bytes"
	"encoding/json"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Helper function to create a test logger
func createTestLogger(level zapcore.Level, buffer *bytes.Buffer) *Logger {
	writer := zapcore.AddSync(buffer)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "" // Remove time for easier testing
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), writer, level)
	return &Logger{Logger: zap.New(core)}
}

// TestSync ensures Sync works without errors.
func TestSync(t *testing.T) {
	logger := GetLogger(false, "test")
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Sync caused a panic: %v", r)
		}
	}()
	logger.Sync() // Should not panic
	t.Log("Sync executed successfully")
}

// TestInfo tests the Info logging method.
func TestInfo(t *testing.T) {
	buffer := &bytes.Buffer{}
	logger := createTestLogger(zap.InfoLevel, buffer)

	// Log an info message
	logger.Info("Info log test", zap.String("key", "value"))

	// Verify the log message
	var loggedData map[string]interface{}
	err := json.Unmarshal(buffer.Bytes(), &loggedData)
	if err != nil {
		t.Fatalf("Failed to parse log message: %v", err)
	}

	if loggedData["msg"] != "Info log test" {
		t.Errorf("Expected message 'Info log test', got '%v'", loggedData["msg"])
	}
	if loggedData["key"] != "value" {
		t.Errorf("Expected field 'key' to be 'value', got '%v'", loggedData["key"])
	}
}

// TestError tests the Error logging method.
func TestError(t *testing.T) {
	buffer := &bytes.Buffer{}
	logger := createTestLogger(zap.ErrorLevel, buffer)

	// Log an error message
	logger.Error("Error log test", zap.String("error", "test_error"))

	// Verify the log message
	var loggedData map[string]interface{}
	err := json.Unmarshal(buffer.Bytes(), &loggedData)
	if err != nil {
		t.Fatalf("Failed to parse log message: %v", err)
	}

	if loggedData["msg"] != "Error log test" {
		t.Errorf("Expected message 'Error log test', got '%v'", loggedData["msg"])
	}
	if loggedData["error"] != "test_error" {
		t.Errorf("Expected field 'error' to be 'test_error', got '%v'", loggedData["error"])
	}
}

// TestDebug tests the Debug logging method.
func TestDebug(t *testing.T) {
	buffer := &bytes.Buffer{}
	logger := createTestLogger(zap.DebugLevel, buffer)

	// Log a debug message
	logger.Debug("Debug log test", zap.Int("line", 42))

	// Verify the log message
	var loggedData map[string]interface{}
	err := json.Unmarshal(buffer.Bytes(), &loggedData)
	if err != nil {
		t.Fatalf("Failed to parse log message: %v", err)
	}

	if loggedData["msg"] != "Debug log test" {
		t.Errorf("Expected message 'Debug log test', got '%v'", loggedData["msg"])
	}
	if loggedData["line"] != float64(42) { // JSON unmarshals numbers to float64
		t.Errorf("Expected field 'line' to be 42, got '%v'", loggedData["line"])
	}
}

// TestWarn tests the Warn logging method.
func TestWarn(t *testing.T) {
	buffer := &bytes.Buffer{}
	logger := createTestLogger(zap.WarnLevel, buffer)

	// Log a warn message
	logger.Warn("Warn log test", zap.String("warning", "disk_space_low"))

	// Verify the log message
	var loggedData map[string]interface{}
	err := json.Unmarshal(buffer.Bytes(), &loggedData)
	if err != nil {
		t.Fatalf("Failed to parse log message: %v", err)
	}

	if loggedData["msg"] != "Warn log test" {
		t.Errorf("Expected message 'Warn log test', got '%v'", loggedData["msg"])
	}
	if loggedData["warning"] != "disk_space_low" {
		t.Errorf("Expected field 'warning' to be 'disk_space_low', got '%v'", loggedData["warning"])
	}
}
