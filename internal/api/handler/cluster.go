package handler

import (
	"context"

	"github.com/tdevsin/keyforge/internal/api/controller"
	"github.com/tdevsin/keyforge/internal/config"
	"github.com/tdevsin/keyforge/internal/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ClusterHandler struct {
	proto.UnimplementedClusterServiceServer
	Conf *config.Config
}

func (c *ClusterHandler) GetClusterState(ctx context.Context, req *emptypb.Empty) (*proto.ClusterState, error) {
	c.Conf.Logger.Info("GetClusterState called")
	return controller.GetClusterInfo(c.Conf)
}

func (c *ClusterHandler) SetClusterState(ctx context.Context, req *proto.ClusterState) (*emptypb.Empty, error) {
	c.Conf.Logger.Info("SetClusterState called")
	return &emptypb.Empty{}, controller.SetClusterInfo(c.Conf, req)
}
