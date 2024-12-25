package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tdevsin/keyforge/internal/proto"
)

func TestSetKey(t *testing.T) {
	_, cleanup := runApp(t)
	defer cleanup()
	conn := getGrpcConnection()
	client := proto.NewKeyServiceClient(conn)

	t.Run("Should set key successfully if all data is valid", func(t *testing.T) {
		request := &proto.SetKeyRequest{
			Key: "k1",
			Value: &proto.SetKeyRequest_StringValue{
				StringValue: "v1",
			},
		}
		response, err := client.SetKey(context.Background(), request)

		assert.Nil(t, err)
		assert.Equal(t, "k1", response.Key)
		assert.Equal(t, "v1", response.GetStringValue())
		assert.Nil(t, response.GetBinaryValue())
	})

	t.Run("Should return binary value if binary value is provided", func(t *testing.T) {
		request := &proto.SetKeyRequest{
			Key: "k2",
			Value: &proto.SetKeyRequest_BinaryValue{
				BinaryValue: []byte{1, 2, 3},
			},
		}
		response, err := client.SetKey(context.Background(), request)

		assert.Nil(t, err)
		assert.Equal(t, "k2", response.Key)
		assert.Nil(t, response.GetStringValue())
		assert.Equal(t, []byte{1, 2, 3}, response.GetBinaryValue())
	})

	t.Run("Should return error if invalid key is provided", func(t *testing.T) {
		request := &proto.SetKeyRequest{
			Value: &proto.SetKeyRequest_StringValue{
				StringValue: "v1",
			},
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
			Key: "k1",
		}
		response, err := client.GetKey(context.Background(), request)

		assert.NotNil(t, err)
		assert.Nil(t, response)
	})

	t.Run("Should return key successfully if key is found", func(t *testing.T) {
		request := &proto.SetKeyRequest{
			Key: "k1",
			Value: &proto.SetKeyRequest_StringValue{
				StringValue: "v1",
			},
		}
		client.SetKey(context.Background(), request)

		getRequest := &proto.GetKeyRequest{
			Key: "k1",
		}
		response, err := client.GetKey(context.Background(), getRequest)

		assert.Nil(t, err)
		assert.Equal(t, "k1", response.Key)
		assert.Equal(t, "v1", response.GetStringValue())
		assert.Nil(t, response.GetBinaryValue())
	})
}

func TestDeleteKey(t *testing.T) {
	_, cleanup := runApp(t)
	defer cleanup()
	conn := getGrpcConnection()
	client := proto.NewKeyServiceClient(conn)

	t.Run("Should return error if key is not found", func(t *testing.T) {
		request := &proto.DeleteKeyRequest{
			Key: "k1",
		}
		response, err := client.DeleteKey(context.Background(), request)

		assert.NotNil(t, err)
		assert.Nil(t, response)
	})

	t.Run("Should delete key successfully if key is found", func(t *testing.T) {
		request := &proto.SetKeyRequest{
			Key: "k1",
			Value: &proto.SetKeyRequest_StringValue{
				StringValue: "v1",
			},
		}
		client.SetKey(context.Background(), request)

		deleteRequest := &proto.DeleteKeyRequest{
			Key: "k1",
		}
		response, err := client.DeleteKey(context.Background(), deleteRequest)

		assert.Nil(t, err)
		assert.Equal(t, "k1", response.Key)

		r := &proto.GetKeyRequest{
			Key: "k1",
		}
		res, err := client.GetKey(context.Background(), r)

		assert.NotNil(t, err)
		assert.Nil(t, res)
	})
}
