package cluster

import (
	"sync"
	"time"
)

type ClusterManager interface {
	GetClusterInfo() *ClusterInfo // Retrieve the current cluster state
	IncrementVersion()            // Increment the cluster state version

	MergeClusterState(receivedState *ClusterInfo) // Merge received cluster state with the current state

	// Node Management
	AddOrUpdateNode(node Node)          // Add or update a node in the cluster
	RemoveNode(nodeID string)           // Remove a node from the cluster
	GetNode(nodeID string) (Node, bool) // Retrieve a node by its ID
	GetHealthyNodes() []Node            // Retrieve a list of healthy nodes

	RegisterObserver(observer ClusterObserver) // Register an observer to get notified on state changes
	// Gossip Propagation
	// PropagateState(toNode Node) error // Propagate the cluster state to a specific node
}

// ClusterInfo represents the overall state of the cluster
type ClusterInfo struct {
	mu          sync.RWMutex      // Mutex to protect concurrent access
	Nodes       map[string]Node   // Nodes is a map of nodeId to Node
	Version     int               // Version helps in identifying the latest cluster state
	LastUpdated time.Time         // LastUpdated indicates the last time the cluster info was updated
	observers   []ClusterObserver // List of observers to notify on state changes
}

func NewCluster() *ClusterInfo {
	return &ClusterInfo{
		Nodes:       make(map[string]Node),
		Version:     -1, // -1 indicates that this node is started for the first time and has no cluster state information
		LastUpdated: time.Now(),
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
	defer ci.mu.Unlock()

	// Ignore outdated states
	if receivedState.Version < ci.Version {
		return
	}

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
		}
	}

	// Update version and metadata if received state is newer
	if receivedState.Version > ci.Version {
		ci.Version = receivedState.Version
	}

	ci.LastUpdated = time.Now()
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
