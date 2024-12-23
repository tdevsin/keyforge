package handler

import (
	"context"

	"github.com/tdevsin/keyforge/internal/logger"
	"github.com/tdevsin/keyforge/internal/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// HealthHandler is the handler for Health Check
type HealthHandler struct {
	proto.UnimplementedHealthServiceServer
}

// CheckHealth returns an empty response indicating the service is healthy
func (*HealthHandler) CheckHealth(context.Context, *emptypb.Empty) (*emptypb.Empty, error) {
	logger.Info("Health Check")
	return &emptypb.Empty{}, nil
}
