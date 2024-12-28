package handler

import (
	"context"

	"github.com/tdevsin/keyforge/internal/config"
	"github.com/tdevsin/keyforge/internal/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// HealthHandler is the handler for Health Check
type HealthHandler struct {
	proto.UnimplementedHealthServiceServer
	Conf *config.Config
}

// CheckHealth returns an empty response indicating the service is healthy
func (h *HealthHandler) CheckHealth(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	h.Conf.Logger.Info("Health Check")
	return &emptypb.Empty{}, nil
}
