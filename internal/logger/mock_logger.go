package logger

import (
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type MockLogging struct {
	mock.Mock
}

func (m *MockLogging) Info(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}

func (m *MockLogging) Error(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}

func (m *MockLogging) Debug(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}

func (m *MockLogging) Warn(msg string, fields ...zap.Field) {
	m.Called(msg, fields)
}

func (m *MockLogging) Sync() {
	m.Called()
}
