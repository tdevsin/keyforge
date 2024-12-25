package test

import (
	"context"
	"testing"

	"github.com/tdevsin/keyforge/internal/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

// TestHealth tests the health endpoint
func TestHealth(t *testing.T) {

	t.Run("Should return success health response if application is running", func(t *testing.T) {
		_, cleanup := runApp(t)
		defer cleanup()
		conn := getGrpcConnection()
		client := proto.NewHealthServiceClient(conn)
		_, err := client.CheckHealth(context.Background(), &emptypb.Empty{})
		if err != nil {
			t.Errorf("Expected %v, Got: %v", "Success", err)
		}
	})

	t.Run("Should return error health response if application is not running", func(t *testing.T) {
		conn := getGrpcConnection()
		client := proto.NewHealthServiceClient(conn)
		_, err := client.CheckHealth(context.Background(), &emptypb.Empty{})
		if err == nil {
			t.Errorf("Expected %v, Got: %v", "Error", err)
		}
	})
}
