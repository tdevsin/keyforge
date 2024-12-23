package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/tdevsin/keyforge/internal/proto"
)

// KVHandler is the handler for Key-Value operations
type KVHandler struct {
	proto.UnimplementedKeyServiceServer
}

// GetKey returns the value for the given key
func (*KVHandler) GetKey(ctx context.Context, req *proto.GetKeyRequest) (*proto.GetKeyResponse, error) {
	return &proto.GetKeyResponse{
		Key: req.Key,
		Value: &proto.GetKeyResponse_StringValue{
			StringValue: uuid.NewString(),
		},
	}, nil
}

// SetKey sets the value for the given key
func (*KVHandler) SetKey(ctx context.Context, req *proto.SetKeyRequest) (*proto.SetKeyResponse, error) {
	return &proto.SetKeyResponse{
		Key: req.Key,
		Value: &proto.SetKeyResponse_StringValue{
			StringValue: req.GetStringValue(),
		},
	}, nil
}

// DeleteKey deletes the key
func (*KVHandler) DeleteKey(ctx context.Context, req *proto.DeleteKeyRequest) (*proto.DeleteKeyResponse, error) {
	return &proto.DeleteKeyResponse{
		Key: req.Key,
	}, nil
}
