package test

// This file contains helper functions that help in running integration tests.

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var appBinary = "./keyforge"

func TestMain(m *testing.M) {
	// Build the application
	err := exec.Command("go", "build", "-o", appBinary, "../main.go").Run()
	if err != nil {
		panic("Failed to build application: " + err.Error())
	}

	// Run the tests
	code := m.Run()

	// Cleanup
	os.Remove(appBinary)

	// Exit with test result code
	os.Exit(code)
}

// runApp start and stop the server for each test
func runApp(t *testing.T) (*exec.Cmd, func()) {
	cmd := exec.Command(appBinary, "start")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Start the server
	err := cmd.Start()
	if err != nil {
		t.Fatalf("Failed to start the application: %v", err)
	}

	// Wait for the server to start
	time.Sleep(2 * time.Second)

	// Cleanup function to stop the server
	cleanup := func() {
		if cmd.Process != nil {
			t.Log("Stopping the application")
			cmd.Process.Kill()
		}
	}

	return cmd, cleanup
}

// getGrpcConnection returns a new GRPC connection with a local instance of the server
func getGrpcConnection() (conn *grpc.ClientConn) {
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		// Panic if connection fails and stop the test
		panic("Failed to connect to server: " + err.Error())
	}
	return conn
}
