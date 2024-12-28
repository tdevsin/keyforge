package storage

import "github.com/stretchr/testify/mock"

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDatabase) WriteKey(key, value []byte) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func (m *MockDatabase) ReadKey(key []byte) ([]byte, error) {
	args := m.Called(key)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockDatabase) DeleteKey(key []byte) error {
	args := m.Called(key)
	return args.Error(0)
}
