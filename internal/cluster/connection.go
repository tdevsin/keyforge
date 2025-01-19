package cluster

import (
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ConnectionPool struct {
	connections map[string]*grpc.ClientConn
	mu          sync.Mutex
}

func NewConnectionPool() *ConnectionPool {
	return &ConnectionPool{connections: make(map[string]*grpc.ClientConn)}
}

func (cp *ConnectionPool) GetConnection(addr string) (*grpc.ClientConn, error) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if conn, exists := cp.connections[addr]; exists {
		return conn, nil
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	cp.connections[addr] = conn
	return conn, nil
}

func (cp *ConnectionPool) Close() {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	for _, conn := range cp.connections {
		conn.Close()
	}
}
