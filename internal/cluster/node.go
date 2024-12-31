package cluster

import "time"

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
	// Unique ID of the Node
	ID string
	// Address of the Node in <host>:<port> format
	Address string
	// Position of this Node on the hash ring
	Position int
	// Health defines health of this Node
	Health Health
}
