package cluster

import (
	"hash/crc32"
	"sort"
)

type ConsistentHashRing interface {
	AddNode(node Node)
	RemoveNode(nodeID string)
	GetResponsibleNode(key string) string
}

type HashRing struct {
	Nodes []Node
}

func NewHashRing() *HashRing {
	return &HashRing{}
}

// AddNode adds a node to the hash ring
func (hr *HashRing) AddNode(node Node) {
	position := CalculateNodePosition(node.ID)
	node.Position = position
	hr.Nodes = append(hr.Nodes, node)
	sort.Slice(hr.Nodes, func(i, j int) bool {
		return hr.Nodes[i].Position < hr.Nodes[j].Position
	})
}

// RemoveNode removes a node from the hash ring
func (hr *HashRing) RemoveNode(nodeID string) {
	position := CalculateNodePosition(nodeID)
	for i, node := range hr.Nodes {
		if node.Position == position {
			hr.Nodes = append(hr.Nodes[:i], hr.Nodes[i+1:]...)
			break
		}
	}
}

// GetResponsibleNode returns the node responsible for a given key
func (hr *HashRing) GetResponsibleNode(key string) string {
	if len(hr.Nodes) == 0 {
		return ""
	}
	keyPosition := CalculateKeyPosition(key)
	for _, node := range hr.Nodes {
		if keyPosition <= node.Position {
			return node.ID
		}
	}
	return hr.Nodes[0].ID // Wrap around to the first node
}

// CalculateNodePosition calculates the position of a node on the ring
func CalculateNodePosition(nodeID string) int {
	hash := crc32.ChecksumIEEE([]byte(nodeID))
	return int(hash)
}

// CalculateKeyPosition calculates the position of a key on the ring
func CalculateKeyPosition(key string) int {
	hash := crc32.ChecksumIEEE([]byte(key))
	return int(hash)
}
