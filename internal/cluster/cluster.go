package cluster

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/tdevsin/keyforge/internal/proto"
	"golang.org/x/exp/rand"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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

// ClusterInfo represents the overall state of the cluster
type ClusterInfo struct {
	mu          sync.RWMutex      // Mutex to protect concurrent access
	Nodes       map[string]Node   // Nodes is a map of nodeId to Node
	Version     int               // Version helps in identifying the latest cluster state
	LastUpdated time.Time         // LastUpdated indicates the last time the cluster info was updated
	observers   []ClusterObserver // List of observers to notify on state changes
	selfId      string            // selfId is the ID of the current node
	gossipN     int               // randomSelection is the number of nodes to select for gossip
}

func NewCluster(selfId string, gossipN int) *ClusterInfo {
	return &ClusterInfo{
		Nodes:       make(map[string]Node),
		Version:     -1, // -1 indicates that this node is started for the first time and has no cluster state information
		LastUpdated: time.Now(),
		selfId:      selfId,
		gossipN:     gossipN,
	}
}

func (ci *ClusterInfo) RegisterObserver(observer ClusterObserver) {
	ci.mu.Lock()
	defer ci.mu.Unlock()
	ci.observers = append(ci.observers, observer)
}

func (ci *ClusterInfo) notifyNodeAdded(node Node) {
	for _, observer := range ci.observers {
		observer.NodeAdded(node)
	}
}

func (ci *ClusterInfo) notifyNodeUpdated(node Node) {
	for _, observer := range ci.observers {
		observer.NodeUpdated(node)
	}
}

func (ci *ClusterInfo) notifyNodeRemoved(nodeID string) {
	for _, observer := range ci.observers {
		observer.NodeRemoved(nodeID)
	}
}

func (ci *ClusterInfo) GetClusterInfo() *ClusterInfo {
	ci.mu.RLock()
	defer ci.mu.RUnlock()
	return ci
}

func (ci *ClusterInfo) IncrementVersion() {
	ci.mu.Lock()
	defer ci.mu.Unlock()
	ci.Version++
	ci.LastUpdated = time.Now()
}

func (ci *ClusterInfo) MergeClusterState(receivedState *ClusterInfo) {
	ci.mu.Lock()

	// Ignore outdated states
	if receivedState.Version < ci.Version {
		ci.mu.Unlock()
		return
	}

	var needToGossip bool
	// Merge node-level data
	for nodeID, receivedNode := range receivedState.Nodes {
		existingNode, exists := ci.Nodes[nodeID]

		// Add new nodes or update existing nodes based on LastChecked
		if !exists || receivedNode.Health.LastChecked.After(existingNode.Health.LastChecked) {
			if !exists {
				ci.notifyNodeAdded(receivedNode)
			} else {
				ci.notifyNodeUpdated(receivedNode)
			}
			ci.Nodes[nodeID] = receivedNode
			needToGossip = true
		}
	}

	// Update version and metadata if received state is newer
	if receivedState.Version > ci.Version {
		ci.Version = receivedState.Version
		needToGossip = true
	}

	ci.LastUpdated = time.Now()

	// Unlock before initiating gossip as gossip also acquires the lock. It will be a deadlock in case it is not unlocked before gossip
	ci.mu.Unlock()
	if needToGossip && ci.Version > -1 {
		ci.InitiateGossip()
	}
}

func (ci *ClusterInfo) AddOrUpdateNode(node Node) {
	ci.mu.Lock()
	defer ci.mu.Unlock()
	if _, exists := ci.Nodes[node.ID]; !exists {
		ci.notifyNodeAdded(node)
	} else {
		ci.notifyNodeUpdated(node)
	}
	ci.Nodes[node.ID] = node
	ci.LastUpdated = time.Now()
}

func (ci *ClusterInfo) RemoveNode(nodeID string) {
	ci.mu.Lock()
	defer ci.mu.Unlock()
	delete(ci.Nodes, nodeID)
	ci.notifyNodeRemoved(nodeID)
	ci.LastUpdated = time.Now()
}

func (ci *ClusterInfo) GetNode(nodeID string) (Node, bool) {
	ci.mu.RLock()
	defer ci.mu.RUnlock()
	node, ok := ci.Nodes[nodeID]
	return node, ok
}

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
func (ci *ClusterInfo) GetRandomNodesForGossip() []Node {
	ci.mu.RLock()
	defer ci.mu.RUnlock()

	// Retrieve all healthy nodes
	healthyNodes := ci.GetHealthyNodes()

	// Filter out the current node (self)
	filteredNodes := make([]Node, 0, len(healthyNodes))
	for _, node := range healthyNodes {
		if node.ID != ci.selfId {
			filteredNodes = append(filteredNodes, node)
		}
	}

	// Determine the effective number of nodes to select
	n := ci.gossipN
	if n > len(filteredNodes) {
		n = len(filteredNodes)
	}

	// Shuffle the remaining nodes
	r := rand.New(rand.NewSource(uint64(time.Now().UnixNano())))
	r.Shuffle(len(filteredNodes), func(i, j int) {
		filteredNodes[i], filteredNodes[j] = filteredNodes[j], filteredNodes[i]
	})

	// Return the first n nodes
	return filteredNodes[:n]
}

func (ci *ClusterInfo) InitiateGossip() error {
	// Get n random nodes
	nodes := ci.GetRandomNodesForGossip()
	var err []error
	// Propagate the cluster state to each of the nodes
	for _, node := range nodes {
		conn, e := grpc.NewClient(node.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if e != nil {
			err = append(err, e)
			continue
		}
		client := proto.NewClusterServiceClient(conn)
		var req proto.ClusterState
		ci.MapClusterStateToProto(&req)
		_, e = client.SetClusterState(context.TODO(), &req)
		if e != nil {
			err = append(err, e)
		}
		conn.Close()
	}
	return errors.Join(err...)
}

func (ci *ClusterInfo) MapClusterStateToProto(state *proto.ClusterState) {
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
