package handler

import (
	"context"

	"github.com/tdevsin/keyforge/internal/api/controller"
	"github.com/tdevsin/keyforge/internal/config"
	"github.com/tdevsin/keyforge/internal/proto"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// KVHandler is the handler for Key-Value operations
type KVHandler struct {
	proto.UnimplementedKeyServiceServer
	Conf *config.Config
}

// GetKey returns the value for the given key
func (k *KVHandler) GetKey(ctx context.Context, req *proto.GetKeyRequest) (*proto.GetKeyResponse, error) {
	k.Conf.Logger.Info("Get Request", zap.Field{Key: req.Key, Type: zapcore.StringType})
	return controller.GetKey(k.Conf, req)
}

// SetKey sets the value for the given key
func (k *KVHandler) SetKey(ctx context.Context, req *proto.SetKeyRequest) (*proto.SetKeyResponse, error) {
	k.Conf.Logger.Info("Set Request", zap.Field{Key: req.Key, Type: zapcore.StringType})
	return controller.SetKey(k.Conf, req)
}

// DeleteKey deletes the key
func (k *KVHandler) DeleteKey(ctx context.Context, req *proto.DeleteKeyRequest) (*proto.DeleteKeyResponse, error) {
	k.Conf.Logger.Info("Delete Request", zap.Field{Key: req.Key, Type: zapcore.StringType})
	return controller.DeleteKey(k.Conf, req)
}
