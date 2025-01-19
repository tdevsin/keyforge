package controller

import (
	"context"

	"github.com/cockroachdb/pebble"
	"github.com/tdevsin/keyforge/internal/config"
	"github.com/tdevsin/keyforge/internal/constants"
	"github.com/tdevsin/keyforge/internal/proto"
	"github.com/tdevsin/keyforge/internal/utils"
	"go.uber.org/zap"
)

func SetKey(ctx context.Context, c *config.Config, r *proto.SetKeyRequest) (*proto.SetKeyResponse, error) {
	if utils.IsEmpty(r.GetKey()) {
		return nil, constants.StatusErrInvalidKey
	}
	if r.GetValue() == nil || len(r.GetKey()) == 0 {
		return nil, constants.StatusErrInvalidValue
	}
	responsibleNode := c.HashRing.GetResponsibleNode(r.GetKey())

	if c.NodeInfo.ID == responsibleNode {
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
	} else {
		return proxySetRequest(ctx, c, c.HashRing.GetNode(responsibleNode).Address, r)
	}
}

func GetKey(ctx context.Context, c *config.Config, r *proto.GetKeyRequest) (*proto.GetKeyResponse, error) {
	if utils.IsEmpty(r.GetKey()) {
		return nil, constants.StatusErrInvalidKey
	}
	responsibleNode := c.HashRing.GetResponsibleNode(r.GetKey())
	if c.NodeInfo.ID == responsibleNode {
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
	} else {
		return proxyGetRequest(ctx, c, c.HashRing.GetNode(responsibleNode).Address, r)
	}
}

func DeleteKey(ctx context.Context, c *config.Config, r *proto.DeleteKeyRequest) (*proto.DeleteKeyResponse, error) {
	if utils.IsEmpty(r.GetKey()) {
		return nil, constants.StatusErrInvalidKey
	}
	responsibleNode := c.HashRing.GetResponsibleNode(r.GetKey())
	if c.NodeInfo.ID == responsibleNode {
		err := c.Db.DeleteKey([]byte(r.GetKey()))
		if err != nil {
			return nil, constants.StatusErrInternal
		}
		return &proto.DeleteKeyResponse{
			Key: r.GetKey(),
		}, nil

	} else {
		return proxyDeleteRequest(ctx, c, c.HashRing.GetNode(responsibleNode).Address, r)
	}
}

func proxyGetRequest(ctx context.Context, conf *config.Config, addr string, request *proto.GetKeyRequest) (*proto.GetKeyResponse, error) {
	conn, err := conf.ConnectionPool.GetConnection(addr)
	if err != nil {
		return nil, err
	}
	client := proto.NewKeyServiceClient(conn)
	return client.GetKey(ctx, request)
}

func proxySetRequest(ctx context.Context, conf *config.Config, addr string, request *proto.SetKeyRequest) (*proto.SetKeyResponse, error) {
	conn, err := conf.ConnectionPool.GetConnection(addr)
	if err != nil {
		return nil, err
	}
	client := proto.NewKeyServiceClient(conn)
	return client.SetKey(ctx, request)
}

func proxyDeleteRequest(ctx context.Context, conf *config.Config, addr string, request *proto.DeleteKeyRequest) (*proto.DeleteKeyResponse, error) {
	conn, err := conf.ConnectionPool.GetConnection(addr)
	if err != nil {
		return nil, err
	}
	client := proto.NewKeyServiceClient(conn)
	return client.DeleteKey(ctx, request)
}
