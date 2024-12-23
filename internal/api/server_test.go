package api

import (
	"testing"
)

func TestGetAvailableListener(t *testing.T) {
	lis, err := getAvailableListener(8080, 8089)
	if err != nil {
		t.Fatalf("Expected an available port, but got error: %v", err)
	}
	defer lis.Close()

	if lis.Addr().String() == "" {
		t.Fatalf("Listener address should not be empty")
	}
}
