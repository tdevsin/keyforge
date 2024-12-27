package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tdevsin/keyforge/internal/config"
	"github.com/tdevsin/keyforge/internal/constants"
	"github.com/tdevsin/keyforge/internal/logger"
	"github.com/tdevsin/keyforge/internal/proto"
	"github.com/tdevsin/keyforge/internal/storage"
)

// Helper function to create a test configuration
func createTestConfig(t *testing.T) *config.Config {
	// Create a temporary directory for the database
	tempDir := t.TempDir()

	// Initialize logger
	log := logger.GetLogger()

	// Initialize database
	db := storage.GetDatabaseInstance(log, tempDir)

	return &config.Config{
		RootDir: tempDir,
		Logger:  log,
		Db:      db,
	}
}

// Test SetKey function
func TestSetKey(t *testing.T) {
	cfg := createTestConfig(t)
	defer cfg.Db.Close()

	t.Run("Invalid Key", func(t *testing.T) {
		req := &proto.SetKeyRequest{
			Key:   "",
			Value: []byte("value"),
		}
		resp, err := SetKey(cfg, req)
		assert.Nil(t, resp)
		assert.Equal(t, constants.StatusErrInvalidKey, err)
	})

	t.Run("Invalid Value", func(t *testing.T) {
		req := &proto.SetKeyRequest{
			Key:   "key1",
			Value: nil,
		}
		resp, err := SetKey(cfg, req)
		assert.Nil(t, resp)
		assert.Equal(t, constants.StatusErrInvalidValue, err)
	})

	t.Run("Successful Write", func(t *testing.T) {
		req := &proto.SetKeyRequest{
			Key:   "key1",
			Value: []byte("value1"),
		}
		resp, err := SetKey(cfg, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "key1", resp.Key)
		assert.Equal(t, []byte("value1"), resp.Value)
	})
}

// Test GetKey function
func TestGetKey(t *testing.T) {
	cfg := createTestConfig(t)
	defer cfg.Db.Close()

	t.Run("Invalid Key", func(t *testing.T) {
		req := &proto.GetKeyRequest{
			Key: "",
		}
		resp, err := GetKey(cfg, req)
		assert.Nil(t, resp)
		assert.Equal(t, constants.StatusErrInvalidKey, err)
	})

	t.Run("Key Not Found", func(t *testing.T) {
		req := &proto.GetKeyRequest{
			Key: "nonexistent-key",
		}
		resp, err := GetKey(cfg, req)
		assert.Nil(t, resp)
		assert.Equal(t, constants.StatusErrKeyNotFound, err)
	})

	t.Run("Successful Read", func(t *testing.T) {
		// Prepopulate the database
		err := cfg.Db.WriteKey([]byte("key1"), []byte("value1"))
		assert.NoError(t, err)

		req := &proto.GetKeyRequest{
			Key: "key1",
		}
		resp, err := GetKey(cfg, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "key1", resp.Key)
		assert.Equal(t, []byte("value1"), resp.Value)
	})
}

// Test DeleteKey function
func TestDeleteKey(t *testing.T) {
	cfg := createTestConfig(t)
	defer cfg.Db.Close()

	t.Run("Invalid Key", func(t *testing.T) {
		req := &proto.DeleteKeyRequest{
			Key: "",
		}
		resp, err := DeleteKey(cfg, req)
		assert.Nil(t, resp)
		assert.Equal(t, constants.StatusErrInvalidKey, err)
	})

	t.Run("Successful Delete", func(t *testing.T) {
		// Prepopulate the database
		err := cfg.Db.WriteKey([]byte("key1"), []byte("value1"))
		assert.NoError(t, err)

		req := &proto.DeleteKeyRequest{
			Key: "key1",
		}
		resp, err := DeleteKey(cfg, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "key1", resp.Key)
	})
}
