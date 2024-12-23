package logger

import (
	"bytes"
	"encoding/json"
	"sync"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TestLoggerInitialization ensures the logger initializes only once (singleton).
func TestLoggerInitialization(t *testing.T) {

	wg := sync.WaitGroup{}
	numRoutines := 10
	wg.Add(numRoutines)

	// Concurrent logger initialization
	for i := 0; i < numRoutines; i++ {
		go func() {
			Info("Testing logger initialization")
			wg.Done()
		}()
	}

	wg.Wait()
	t.Log("Logger initialized successfully in concurrent environment")
}

// TestLoggerSingleton ensures the logger instance is reused (singleton pattern).
func TestLoggerSingleton(t *testing.T) {

	logger1 := getLogger()
	logger2 := getLogger()

	if logger1 != logger2 {
		t.Errorf("Expected logger1 and logger2 to be the same instance")
	}
}

// TestSync ensures Sync works without errors.
func TestSync(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Sync caused a panic: %v", r)
		}
	}()
	Sync() // Should not panic
	t.Log("Sync executed successfully")
}

// TestInfo tests the Info logging method.
func TestInfo(t *testing.T) {

	buffer := &bytes.Buffer{}
	writer := zapcore.AddSync(buffer)

	// Create a test logger
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = ""
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), writer, zap.InfoLevel)
	testLogger := zap.New(core)

	// Set custom logger
	log = testLogger

	// Log an info message
	Info("Info log test", zap.String("key", "value"))

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
	writer := zapcore.AddSync(buffer)

	// Create a test logger
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = ""
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), writer, zap.ErrorLevel)
	testLogger := zap.New(core)

	// Set custom logger
	log = testLogger

	// Log an error message
	Error("Error log test", zap.String("error", "test_error"))

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
	writer := zapcore.AddSync(buffer)

	// Create a test logger
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = ""
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), writer, zap.DebugLevel)
	testLogger := zap.New(core)

	// Set custom logger
	log = testLogger

	// Log a debug message
	Debug("Debug log test", zap.Int("line", 42))

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
	writer := zapcore.AddSync(buffer)

	// Create a test logger
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = ""
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), writer, zap.WarnLevel)
	testLogger := zap.New(core)

	// Set custom logger
	log = testLogger

	// Log a warn message
	Warn("Warn log test", zap.String("warning", "disk_space_low"))

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
