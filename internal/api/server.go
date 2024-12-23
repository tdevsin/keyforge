package api

import (
	"fmt"
	"net"

	"github.com/google/uuid"
	"github.com/tdevsin/keyforge/internal/api/handler"
	"github.com/tdevsin/keyforge/internal/logger"
	"github.com/tdevsin/keyforge/internal/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// getAvailableListener tries to find an available port within the given range
func getAvailableListener(start, end int) (net.Listener, error) {
	for i := start; i <= end; i++ {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", i))
		if err == nil {
			return lis, nil
		}
	}
	return nil, fmt.Errorf("no available ports in range %d-%d", start, end)
}

// StartGRPCServer starts a GRPC server on a port between 8080 and 8100
func StartGRPCServer() error {
	lis, err := getAvailableListener(8080, 8100)
	if err != nil {
		logger.Error("Failed to find an available port", zap.Error(err))
		return err
	}

	// Setup gRPC server
	server := grpc.NewServer()

	// Reflection is used by clients like Postman to list services on the server and understand what methods are available
	reflection.Register(server)

	// Generate a new NodeId. This will be used to uniquely identify the node
	nodeId := uuid.New().String()
	logger.Info("Starting GRPC Server", zap.String("port", lis.Addr().String()), zap.String("nodeId", nodeId))

	// Register services
	proto.RegisterKeyServiceServer(server, &handler.KVHandler{})
	proto.RegisterHealthServiceServer(server, &handler.HealthHandler{})

	// Serve the server
	if err := server.Serve(lis); err != nil {
		logger.Error("Failed to start GRPC Server", zap.Error(err))
		return err
	}
	return nil
}
