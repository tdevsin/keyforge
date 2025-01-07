package cluster

import (
	"time"

	"google.golang.org/grpc"
)

type Status int

const (
	Healthy Status = iota
	SuspectedFailed
	PermanentFailed
)

type Health struct {
	Status      Status
	LastChecked time.Time
}

// Node defines a single node in the cluster
type Node struct {
	ID          string           // Unique ID of the Node
	Address     string           // Address of the Node in <host>:<port> format
	Position    int              // Position of this Node on the hash ring
	Health      Health           // Health defines health of this Node
	HealthConn  *grpc.ClientConn // HealthConn is a gRPC connection to the health server of this Node
	ClusterConn *grpc.ClientConn // ClusterConn is a gRPC connection to the cluster server of this Node
}
