package cluster

// ClusterObserver is an interface that allows listening to cluster state changes and updates
// Any struct can implement this interface and register in cluster state to get notified on state changes
type ClusterObserver interface {
	NodeAdded(node Node)
	NodeRemoved(nodeID string)
	NodeHealthSuspectedFailed(nodeId string)
	NodeHealthPermanentFailed(nodeId string)
}
