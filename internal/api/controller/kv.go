package controller

import (
	"github.com/tdevsin/keyforge/internal/constants"
	"github.com/tdevsin/keyforge/internal/logger"
	"github.com/tdevsin/keyforge/internal/proto"
	"github.com/tdevsin/keyforge/internal/storage"
	"github.com/tdevsin/keyforge/internal/utils"
	"go.uber.org/zap"
)

func SetKey(r *proto.SetKeyRequest) (*proto.SetKeyResponse, error) {
	if utils.IsEmpty(r.GetKey()) {
		return nil, constants.ErrInvalidKey
	}
	if r.GetValue() == nil || len(r.GetKey()) == 0 {
		return nil, constants.ErrInvalidValue
	}
	db := storage.GetDatabaseInstance()
	err := db.WriteKey([]byte(r.GetKey()), r.GetValue())
	if err != nil {
		logger.Error("Some error occurred while writing key", zap.Error(err))
		return nil, constants.ErrInternal
	}
	// Set key in DB
	return &proto.SetKeyResponse{
		Key:   r.GetKey(),
		Value: r.GetValue(),
	}, nil
}

func GetKey(r *proto.GetKeyRequest) (*proto.GetKeyResponse, error) {
	if utils.IsEmpty(r.GetKey()) {
		return nil, constants.ErrInvalidKey
	}

	// Get key from db
	db := storage.GetDatabaseInstance()
	v, err := db.ReadKey([]byte(r.GetKey()))
	if err != nil {
		return nil, constants.ErrKeyNotFound
	}
	return &proto.GetKeyResponse{
		Key:   r.GetKey(),
		Value: v,
	}, nil
}
