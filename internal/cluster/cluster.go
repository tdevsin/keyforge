// Package cluster manages the cluster state and provides utilities for interacting with it.
package cluster

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/tdevsin/keyforge/internal/proto"
	"golang.org/x/exp/rand"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ClusterManager defines the interface for cluster operations.
type ClusterManager interface {
	GetClusterInfo() *ClusterInfo                     // Retrieve the current cluster state
	IncrementVersion()                                // Increment the cluster state version
	MergeClusterState(receivedState *ClusterInfo)     // Merge received cluster state with the current state
	AddOrUpdateNode(node Node)                        // Add or update a node in the cluster
	RemoveNode(nodeID string)                         // Remove a node from the cluster
	GetNode(nodeID string) (Node, bool)               // Retrieve a node by its ID
	GetHealthyNodes() []Node                          // Retrieve a list of healthy nodes
	RegisterObserver(observer ClusterObserver)        // Register an observer to get notified on state changes
	MapClusterStateToProto(state *proto.ClusterState) // Map the cluster state to proto
}

// ClusterInfo represents the overall state of the cluster.
type ClusterInfo struct {
	mu             sync.RWMutex      // Mutex to protect concurrent access
	Nodes          map[string]Node   // Nodes is a map of nodeId to Node
	Version        int               // Version helps in identifying the latest cluster state
	LastUpdated    time.Time         // LastUpdated indicates the last time the cluster info was updated
	observers      []ClusterObserver // List of observers to notify on state changes
	selfId         string            // selfId is the ID of the current node
	gossipN        int               // Number of nodes to select for gossip
	gossipInterval time.Duration     // Duration between which the cluster state sync will happen
}

// NewCluster creates and initializes a new ClusterInfo.
func NewCluster(selfId string, gossipN int) *ClusterInfo {
	cluster := &ClusterInfo{
		Nodes:          make(map[string]Node),
		Version:        -1, // Indicates the node is starting for the first time
		LastUpdated:    time.Now(),
		selfId:         selfId,
		gossipN:        gossipN,
		gossipInterval: time.Second * 10,
	}

	return cluster
}

// RegisterObserver registers a new observer for cluster state changes.
func (ci *ClusterInfo) RegisterObserver(observer ClusterObserver) {
	ci.mu.Lock()
	defer ci.mu.Unlock()
	ci.observers = append(ci.observers, observer)
}

// notifyObservers notifies all registered observers of a specific event.
func (ci *ClusterInfo) notifyObservers(event string, nodeID string, node *Node) {
	ci.mu.RLock()
	observers := append([]ClusterObserver(nil), ci.observers...)
	ci.mu.RUnlock()

	for _, observer := range observers {
		switch event {
		case "added":
			if node != nil {
				observer.NodeAdded(*node)
			}
		case "updated":
			if node != nil {
				observer.NodeUpdated(*node)
			}
		case "removed":
			observer.NodeRemoved(nodeID)
		}
	}
}

// GetClusterInfo provides a snapshot of the current cluster state.
func (ci *ClusterInfo) GetClusterInfo() *ClusterInfo {
	ci.mu.RLock()
	defer ci.mu.RUnlock()
	return ci
}

// IncrementVersion increments the cluster state version.
func (ci *ClusterInfo) IncrementVersion() {
	ci.mu.Lock()
	defer ci.mu.Unlock()
	ci.Version++
	ci.LastUpdated = time.Now()
}

// MergeClusterState merges a received cluster state with the current one.
func (ci *ClusterInfo) MergeClusterState(receivedState *ClusterInfo) {
	ci.mu.Lock()

	if receivedState.Version < ci.Version {
		ci.mu.Unlock()
		return
	}

	var addedNodes []Node
	var updatedNodes []Node

	for nodeID, receivedNode := range receivedState.Nodes {
		existingNode, exists := ci.Nodes[nodeID]

		if !exists || receivedNode.Health.LastChecked.After(existingNode.Health.LastChecked) {
			ci.Nodes[nodeID] = receivedNode
			if !exists {
				addedNodes = append(addedNodes, receivedNode)
			} else {
				updatedNodes = append(updatedNodes, receivedNode)
			}
		}
	}

	if receivedState.Version >= ci.Version {
		ci.Version = receivedState.Version
	}
	ci.LastUpdated = time.Now()
	ci.mu.Unlock()

	for _, node := range addedNodes {
		ci.notifyObservers("added", node.ID, &node)
	}
	for _, node := range updatedNodes {
		ci.notifyObservers("updated", node.ID, &node)
	}
}

// AddOrUpdateNode adds or updates a node in the cluster.
func (ci *ClusterInfo) AddOrUpdateNode(node Node) {
	ci.mu.Lock()
	isNewNode := false

	if _, exists := ci.Nodes[node.ID]; !exists {
		isNewNode = true
	}
	ci.Nodes[node.ID] = node
	ci.LastUpdated = time.Now()
	ci.mu.Unlock()

	if isNewNode {
		ci.notifyObservers("added", node.ID, &node)
	} else {
		ci.notifyObservers("updated", node.ID, &node)
	}
}

// RemoveNode removes a node from the cluster.
func (ci *ClusterInfo) RemoveNode(nodeID string) {
	ci.mu.Lock()
	delete(ci.Nodes, nodeID)
	ci.LastUpdated = time.Now()
	ci.mu.Unlock()

	ci.notifyObservers("removed", nodeID, nil)
}

// GetNode retrieves a node by its ID.
func (ci *ClusterInfo) GetNode(nodeID string) (Node, bool) {
	ci.mu.RLock()
	defer ci.mu.RUnlock()
	node, exists := ci.Nodes[nodeID]
	return node, exists
}

// GetHealthyNodes retrieves all healthy nodes.
func (ci *ClusterInfo) GetHealthyNodes() []Node {
	ci.mu.RLock()
	defer ci.mu.RUnlock()
	var healthyNodes []Node
	for _, node := range ci.Nodes {
		if node.Health.Status == Healthy {
			healthyNodes = append(healthyNodes, node)
		}
	}
	return healthyNodes
}

// GetRandomNodesForGossip selects random nodes for gossip.
func (ci *ClusterInfo) GetRandomNodesForGossip() []Node {
	healthyNodes := ci.GetHealthyNodes()

	filteredNodes := make([]Node, 0, len(healthyNodes))
	for _, node := range healthyNodes {
		if node.ID != ci.selfId {
			filteredNodes = append(filteredNodes, node)
		}
	}

	n := ci.gossipN
	if n > len(filteredNodes) {
		n = len(filteredNodes)
	}

	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	r.Shuffle(len(filteredNodes), func(i, j int) {
		filteredNodes[i], filteredNodes[j] = filteredNodes[j], filteredNodes[i]
	})

	return filteredNodes[:n]
}

// createClusterSnapshot creates a consistent snapshot of the current cluster state.
func (ci *ClusterInfo) createClusterSnapshot() *ClusterInfo {
	ci.mu.RLock()
	defer ci.mu.RUnlock()

	copiedNodes := make(map[string]Node, len(ci.Nodes))
	for id, node := range ci.Nodes {
		copiedNodes[id] = node
	}

	return &ClusterInfo{
		Nodes:       copiedNodes,
		Version:     ci.Version,
		LastUpdated: ci.LastUpdated,
		selfId:      ci.selfId,
		gossipN:     ci.gossipN,
	}
}

// InitiateGossip propagates the cluster state to random nodes.
func (ci *ClusterInfo) InitiateGossip() error {
	nodes := ci.GetRandomNodesForGossip()
	var errs []error

	for _, node := range nodes {
		conn, err := grpc.NewClient(node.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			errs = append(errs, err)
			continue
		}
		defer conn.Close()

		client := proto.NewClusterServiceClient(conn)
		var req proto.ClusterState
		ci.MapClusterStateToProto(&req)
		_, err = client.SetClusterState(context.TODO(), &req)
		if err != nil {
			errs = append(errs, err)
		}
		fmt.Printf("**********\nSending state with version %v to Node %v\n****************\n", ci.Version, node.Address)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// MapClusterStateToProto maps the current cluster state to a proto message.
func (ci *ClusterInfo) MapClusterStateToProto(state *proto.ClusterState) {
	ci.mu.RLock()
	defer ci.mu.RUnlock()

	state.Version = int64(ci.Version)
	state.LastUpdated = timestamppb.New(ci.LastUpdated)
	state.Nodes = make([]*proto.Node, 0, len(ci.Nodes))
	for _, node := range ci.Nodes {
		state.Nodes = append(state.Nodes, &proto.Node{
			Id:      node.ID,
			Address: node.Address,
			Health: &proto.Health{
				LastUpdated: timestamppb.New(node.Health.LastChecked),
				Status:      proto.Status(node.Health.Status),
			},
		})
	}
}

// ClusterInfo implements the ClusterObserver interface.
func (ci *ClusterInfo) NodeAdded(node Node) {
	// Trigger gossip when a new node is added
	ci.startGossip()
}

func (ci *ClusterInfo) NodeUpdated(node Node) {
	// Trigger gossip when a node is updated
	ci.startGossip()
}

func (ci *ClusterInfo) NodeRemoved(nodeID string) {
	// Trigger gossip when a node is removed
	ci.startGossip()
}

// startGossip handles initiating gossip in a separate goroutine.
func (ci *ClusterInfo) startGossip() {
	// Create a consistent snapshot of the cluster state for gossip
	snapshot := ci.createClusterSnapshot()
	go func() {
		if err := snapshot.InitiateGossip(); err != nil {
			fmt.Printf("Error during gossip: %v\n", err)
		}
	}()
}

func (ci *ClusterInfo) StartPeriodicGossip() {
	go func() {
		ticker := time.NewTicker(ci.gossipInterval)
		defer ticker.Stop()

		for range ticker.C {
			fmt.Println("Initiating periodic gossip...")
			ci.startGossip()
		}
	}()
}
