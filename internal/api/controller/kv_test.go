package controller

import (
	"context"
	"errors"
	"testing"

	"github.com/cockroachdb/pebble"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tdevsin/keyforge/internal/cluster"
	"github.com/tdevsin/keyforge/internal/config"
	"github.com/tdevsin/keyforge/internal/constants"
	"github.com/tdevsin/keyforge/internal/logger"
	"github.com/tdevsin/keyforge/internal/proto"
	"github.com/tdevsin/keyforge/internal/storage"
)

func TestSetKey(t *testing.T) {
	t.Run("Invalid Key", func(t *testing.T) {
		mockDb := new(storage.MockDatabase)
		mockLogger := new(logger.MockLogging)

		c := &config.Config{
			Db:     mockDb,
			Logger: mockLogger,
		}

		req := &proto.SetKeyRequest{
			Key:   "",
			Value: []byte("value"),
		}

		resp, err := SetKey(context.TODO(), c, req)

		assert.Nil(t, resp)
		assert.Equal(t, constants.StatusErrInvalidKey, err)
	})

	t.Run("Invalid Value", func(t *testing.T) {
		mockDb := new(storage.MockDatabase)
		mockLogger := new(logger.MockLogging)

		c := &config.Config{
			Db:     mockDb,
			Logger: mockLogger,
		}

		req := &proto.SetKeyRequest{
			Key:   "key",
			Value: nil,
		}

		resp, err := SetKey(context.TODO(), c, req)

		assert.Nil(t, resp)
		assert.Equal(t, constants.StatusErrInvalidValue, err)
	})

	t.Run("Database Write Error", func(t *testing.T) {
		mockDb := new(storage.MockDatabase)
		mockLogger := new(logger.MockLogging)

		mockDb.On("WriteKey", []byte("key"), []byte("value")).Return(pebble.ErrReadOnly)
		mockLogger.On("Error", "Some error occurred while writing key", mock.Anything)
		id := uuid.NewString()
		node := cluster.Node{
			ID: id,
		}
		hashring := cluster.NewHashRing()
		hashring.AddNode(node)
		c := &config.Config{
			Db:       mockDb,
			Logger:   mockLogger,
			NodeInfo: &node,
			HashRing: hashring,
		}

		req := &proto.SetKeyRequest{
			Key:   "key",
			Value: []byte("value"),
		}

		resp, err := SetKey(context.TODO(), c, req)

		assert.Nil(t, resp)
		assert.Equal(t, constants.StatusErrInternal, err)

		mockDb.AssertExpectations(t)
		mockLogger.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		mockDb := new(storage.MockDatabase)
		mockLogger := new(logger.MockLogging)
		id := uuid.NewString()
		node := cluster.Node{
			ID: id,
		}
		hashring := cluster.NewHashRing()
		hashring.AddNode(node)
		mockDb.On("WriteKey", []byte("key"), []byte("value")).Return(nil)

		c := &config.Config{
			Db:       mockDb,
			Logger:   mockLogger,
			NodeInfo: &node,
			HashRing: hashring,
		}

		req := &proto.SetKeyRequest{
			Key:   "key",
			Value: []byte("value"),
		}

		resp, err := SetKey(context.TODO(), c, req)

		assert.NotNil(t, resp)
		assert.Equal(t, &proto.SetKeyResponse{
			Key:   "key",
			Value: []byte("value"),
		}, resp)
		assert.Nil(t, err)

		mockDb.AssertExpectations(t)
	})
}

func TestGetKey(t *testing.T) {
	t.Run("Invalid Key", func(t *testing.T) {
		mockDb := new(storage.MockDatabase)
		mockLogger := new(logger.MockLogging)

		c := &config.Config{
			Db:     mockDb,
			Logger: mockLogger,
		}

		req := &proto.GetKeyRequest{
			Key: "",
		}

		resp, err := GetKey(context.TODO(), c, req)

		assert.Nil(t, resp)
		assert.Equal(t, constants.StatusErrInvalidKey, err)
	})

	t.Run("Key Not Found", func(t *testing.T) {
		mockDb := new(storage.MockDatabase)
		mockLogger := new(logger.MockLogging)

		mockDb.On("ReadKey", []byte("key")).Return([]byte(nil), pebble.ErrNotFound)
		id := uuid.NewString()
		node := cluster.Node{
			ID: id,
		}
		hashring := cluster.NewHashRing()
		hashring.AddNode(node)
		c := &config.Config{
			Db:       mockDb,
			Logger:   mockLogger,
			NodeInfo: &node,
			HashRing: hashring,
		}

		req := &proto.GetKeyRequest{
			Key: "key",
		}

		resp, err := GetKey(context.TODO(), c, req)

		assert.Nil(t, resp)
		assert.Equal(t, constants.StatusErrKeyNotFound, err)

		mockDb.AssertExpectations(t)
	})

	t.Run("Database Error", func(t *testing.T) {
		mockDb := new(storage.MockDatabase)
		mockLogger := new(logger.MockLogging)

		mockDb.On("ReadKey", []byte("key")).Return([]byte(nil), errors.New("db error"))
		id := uuid.NewString()
		node := cluster.Node{
			ID: id,
		}
		hashring := cluster.NewHashRing()
		hashring.AddNode(node)
		c := &config.Config{
			Db:       mockDb,
			Logger:   mockLogger,
			NodeInfo: &node,
			HashRing: hashring,
		}
		req := &proto.GetKeyRequest{
			Key: "key",
		}

		resp, err := GetKey(context.TODO(), c, req)

		assert.Nil(t, resp)
		assert.Equal(t, constants.StatusErrInternal, err)

		mockDb.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		mockDb := new(storage.MockDatabase)
		mockLogger := new(logger.MockLogging)
		id := uuid.NewString()
		node := cluster.Node{
			ID: id,
		}
		hashring := cluster.NewHashRing()
		hashring.AddNode(node)
		mockDb.On("ReadKey", []byte("key")).Return([]byte("value"), nil)

		c := &config.Config{
			Db:       mockDb,
			Logger:   mockLogger,
			NodeInfo: &node,
			HashRing: hashring,
		}

		req := &proto.GetKeyRequest{
			Key: "key",
		}

		resp, err := GetKey(context.TODO(), c, req)

		assert.NotNil(t, resp)
		assert.Equal(t, &proto.GetKeyResponse{
			Key:   "key",
			Value: []byte("value"),
		}, resp)
		assert.Nil(t, err)

		mockDb.AssertExpectations(t)
	})
}

func TestDeleteKey(t *testing.T) {
	t.Run("Invalid Key", func(t *testing.T) {
		mockDb := new(storage.MockDatabase)
		mockLogger := new(logger.MockLogging)

		c := &config.Config{
			Db:     mockDb,
			Logger: mockLogger,
		}

		req := &proto.DeleteKeyRequest{
			Key: "",
		}

		resp, err := DeleteKey(context.TODO(), c, req)

		assert.Nil(t, resp)
		assert.Equal(t, constants.StatusErrInvalidKey, err)
	})

	t.Run("Database Error", func(t *testing.T) {
		mockDb := new(storage.MockDatabase)
		mockLogger := new(logger.MockLogging)

		mockDb.On("DeleteKey", []byte("key")).Return(errors.New("db error"))
		id := uuid.NewString()
		node := cluster.Node{
			ID: id,
		}
		hashring := cluster.NewHashRing()
		hashring.AddNode(node)
		c := &config.Config{
			Db:       mockDb,
			Logger:   mockLogger,
			NodeInfo: &node,
			HashRing: hashring,
		}

		req := &proto.DeleteKeyRequest{
			Key: "key",
		}

		resp, err := DeleteKey(context.TODO(), c, req)

		assert.Nil(t, resp)
		assert.Equal(t, constants.StatusErrInternal, err)

		mockDb.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		mockDb := new(storage.MockDatabase)
		mockLogger := new(logger.MockLogging)
		id := uuid.NewString()
		node := cluster.Node{
			ID: id,
		}
		hashring := cluster.NewHashRing()
		hashring.AddNode(node)
		mockDb.On("DeleteKey", []byte("key")).Return(nil)

		c := &config.Config{
			Db:       mockDb,
			Logger:   mockLogger,
			NodeInfo: &node,
			HashRing: hashring,
		}

		req := &proto.DeleteKeyRequest{
			Key: "key",
		}

		resp, err := DeleteKey(context.TODO(), c, req)

		assert.NotNil(t, resp)
		assert.Equal(t, &proto.DeleteKeyResponse{
			Key: "key",
		}, resp)
		assert.Nil(t, err)

		mockDb.AssertExpectations(t)
	})
}
