package handler

import (
	"context"

	"github.com/tdevsin/keyforge/internal/api/controller"
	"github.com/tdevsin/keyforge/internal/proto"
)

// KVHandler is the handler for Key-Value operations
type KVHandler struct {
	proto.UnimplementedKeyServiceServer
}

// GetKey returns the value for the given key
func (*KVHandler) GetKey(ctx context.Context, req *proto.GetKeyRequest) (*proto.GetKeyResponse, error) {
	return controller.GetKey(req)
}

// SetKey sets the value for the given key
func (*KVHandler) SetKey(ctx context.Context, req *proto.SetKeyRequest) (*proto.SetKeyResponse, error) {
	return controller.SetKey(req)
}

// DeleteKey deletes the key
func (*KVHandler) DeleteKey(ctx context.Context, req *proto.DeleteKeyRequest) (*proto.DeleteKeyResponse, error) {
	return &proto.DeleteKeyResponse{
		Key: req.Key,
	}, nil
}
