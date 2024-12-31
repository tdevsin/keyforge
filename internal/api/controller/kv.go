package controller

import (
	"github.com/cockroachdb/pebble"
	"github.com/tdevsin/keyforge/internal/config"
	"github.com/tdevsin/keyforge/internal/constants"
	"github.com/tdevsin/keyforge/internal/proto"
	"github.com/tdevsin/keyforge/internal/utils"
	"go.uber.org/zap"
)

func SetKey(c *config.Config, r *proto.SetKeyRequest) (*proto.SetKeyResponse, error) {
	if utils.IsEmpty(r.GetKey()) {
		return nil, constants.StatusErrInvalidKey
	}
	if r.GetValue() == nil || len(r.GetKey()) == 0 {
		return nil, constants.StatusErrInvalidValue
	}
	if c.NodeInfo.ID == c.HashRing.GetResponsibleNode(r.GetKey()) {
		err := c.Db.WriteKey([]byte(r.GetKey()), r.GetValue())
		if err != nil {
			c.Logger.Error("Some error occurred while writing key", zap.Error(err))
			return nil, constants.StatusErrInternal
		}
		// Set key in DB
		return &proto.SetKeyResponse{
			Key:   r.GetKey(),
			Value: r.GetValue(),
		}, nil
	}
	return &proto.SetKeyResponse{}, nil
}

func GetKey(c *config.Config, r *proto.GetKeyRequest) (*proto.GetKeyResponse, error) {
	if utils.IsEmpty(r.GetKey()) {
		return nil, constants.StatusErrInvalidKey
	}

	if c.NodeInfo.ID == c.HashRing.GetResponsibleNode(r.GetKey()) {
		// Get key from db
		v, err := c.Db.ReadKey([]byte(r.GetKey()))
		if err != nil {
			if err == pebble.ErrNotFound {
				return nil, constants.StatusErrKeyNotFound
			} else {
				return nil, constants.StatusErrInternal
			}
		}
		return &proto.GetKeyResponse{
			Key:   r.GetKey(),
			Value: v,
		}, nil
	}
	return &proto.GetKeyResponse{}, nil
}

func DeleteKey(c *config.Config, r *proto.DeleteKeyRequest) (*proto.DeleteKeyResponse, error) {
	if utils.IsEmpty(r.GetKey()) {
		return nil, constants.StatusErrInvalidKey
	}
	if c.NodeInfo.ID == c.HashRing.GetResponsibleNode(r.GetKey()) {
		err := c.Db.DeleteKey([]byte(r.GetKey()))
		if err != nil {
			return nil, constants.StatusErrInternal
		}
		return &proto.DeleteKeyResponse{
			Key: r.GetKey(),
		}, nil

	}
	return &proto.DeleteKeyResponse{}, nil
}
