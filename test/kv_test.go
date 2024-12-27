package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tdevsin/keyforge/internal/proto"
)

func TestSetKey(t *testing.T) {
	_, cleanup := runApp(t)
	defer cleanup()
	conn := getGrpcConnection()
	client := proto.NewKeyServiceClient(conn)

	t.Run("Should set key successfully if all data is valid", func(t *testing.T) {
		key := "k1" + time.Now().String()
		request := &proto.SetKeyRequest{
			Key:   key,
			Value: []byte("v1"),
		}
		response, err := client.SetKey(context.Background(), request)

		assert.Nil(t, err)
		assert.Equal(t, key, response.GetKey())
		assert.Equal(t, "v1", string(response.GetValue()))
	})

	t.Run("Should return error if invalid key is provided", func(t *testing.T) {
		request := &proto.SetKeyRequest{
			Value: []byte("v1"),
		}
		response, err := client.SetKey(context.Background(), request)

		assert.NotNil(t, err)
		assert.Nil(t, response)
	})
}

func TestGetKey(t *testing.T) {
	_, cleanup := runApp(t)
	defer cleanup()
	conn := getGrpcConnection()
	client := proto.NewKeyServiceClient(conn)

	t.Run("Should return error if key is not found", func(t *testing.T) {
		request := &proto.GetKeyRequest{
			Key: "k1" + time.Now().String(),
		}
		response, err := client.GetKey(context.Background(), request)

		assert.NotNil(t, err)
		assert.Nil(t, response)
	})

	t.Run("Should return key successfully if key is found", func(t *testing.T) {
		key := "k1" + time.Now().String()
		request := &proto.SetKeyRequest{
			Key:   key,
			Value: []byte("v1"),
		}
		client.SetKey(context.Background(), request)

		getRequest := &proto.GetKeyRequest{
			Key: key,
		}
		response, err := client.GetKey(context.Background(), getRequest)

		assert.Nil(t, err)
		assert.Equal(t, key, response.GetKey())
		assert.Equal(t, "v1", string(response.GetValue()))
	})
}

func TestDeleteKey(t *testing.T) {
	_, cleanup := runApp(t)
	defer cleanup()
	conn := getGrpcConnection()
	client := proto.NewKeyServiceClient(conn)

	t.Run("Should delete key successfully", func(t *testing.T) {
		key := "k1" + time.Now().String()
		request := &proto.SetKeyRequest{
			Key:   key,
			Value: []byte("v1"),
		}
		client.SetKey(context.Background(), request)

		deleteRequest := &proto.DeleteKeyRequest{
			Key: key,
		}
		response, err := client.DeleteKey(context.Background(), deleteRequest)

		assert.Nil(t, err)
		assert.Equal(t, key, response.Key)

		r := &proto.GetKeyRequest{
			Key: key,
		}
		res, err := client.GetKey(context.Background(), r)

		assert.NotNil(t, err)
		assert.Nil(t, res)
	})
}
