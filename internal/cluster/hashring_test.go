package cluster

import (
	"hash/crc32"
	"testing"
)

func TestHashRing(t *testing.T) {
	ring := NewHashRing()

	t.Run("AddNodes", func(t *testing.T) {
		nodes := []Node{
			{ID: "NodeA", Address: "127.0.0.1:8000"},
			{ID: "NodeB", Address: "127.0.0.1:8001"},
			{ID: "NodeC", Address: "127.0.0.1:8002"},
		}

		for _, node := range nodes {
			ring.AddNode(node)
		}

		// Ensure nodes are sorted correctly by position
		for i := 1; i < len(ring.Nodes); i++ {
			if ring.Nodes[i].Position < ring.Nodes[i-1].Position {
				t.Fatalf("Nodes are not sorted by position")
			}
		}
	})

	t.Run("KeyToNodeMapping", func(t *testing.T) {
		// Explicitly calculate key positions and check responsible nodes
		keys := map[string]string{}
		for _, key := range []string{"key1", "key2", "key3"} {
			keyPosition := CalculateKeyPosition(key)
			expectedNode := ""
			for _, node := range ring.Nodes {
				if keyPosition <= node.Position {
					expectedNode = node.ID
					break
				}
			}
			if expectedNode == "" {
				expectedNode = ring.Nodes[0].ID // Wrap around to the first node
			}
			keys[key] = expectedNode
		}

		// Validate each key maps to the expected node
		for key, expectedNode := range keys {
			responsibleNode := ring.GetResponsibleNode(key)
			if responsibleNode != expectedNode {
				t.Errorf("For key '%s', expected node '%s', but got '%s'", key, expectedNode, responsibleNode)
			}
		}
	})

	t.Run("RemoveNode", func(t *testing.T) {
		// Remove NodeB
		ring.RemoveNode("NodeB")

		// Ensure NodeB is no longer in the ring
		for _, node := range ring.Nodes {
			if node.ID == "NodeB" {
				t.Fatalf("NodeB should have been removed, but it still exists")
			}
		}

		// Validate key mappings after removal
		keys := map[string]string{}
		for _, key := range []string{"key1", "key2", "key3"} {
			keyPosition := CalculateKeyPosition(key)
			expectedNode := ""
			for _, node := range ring.Nodes {
				if keyPosition <= node.Position {
					expectedNode = node.ID
					break
				}
			}
			if expectedNode == "" {
				expectedNode = ring.Nodes[0].ID // Wrap around to the first node
			}
			keys[key] = expectedNode
		}

		// Validate each key maps to the expected node after node removal
		for key, expectedNode := range keys {
			responsibleNode := ring.GetResponsibleNode(key)
			if responsibleNode != expectedNode {
				t.Errorf("After removal, for key '%s', expected node '%s', but got '%s'", key, expectedNode, responsibleNode)
			}
		}
	})

	t.Run("EmptyRing", func(t *testing.T) {
		emptyRing := NewHashRing()

		// Test GetResponsibleNode on an empty ring
		key := "key1"
		responsibleNode := emptyRing.GetResponsibleNode(key)
		if responsibleNode != "" {
			t.Errorf("For an empty ring, expected no responsible node, but got '%s'", responsibleNode)
		}
	})
}

func TestHashCalculations(t *testing.T) {
	t.Run("NodeHashCalculation", func(t *testing.T) {
		// Test hashing for nodes
		nodeIDs := []string{"NodeA", "NodeB", "NodeC"}
		expectedHashes := []int{
			int(crc32.ChecksumIEEE([]byte("NodeA"))),
			int(crc32.ChecksumIEEE([]byte("NodeB"))),
			int(crc32.ChecksumIEEE([]byte("NodeC"))),
		}

		for i, nodeID := range nodeIDs {
			calculatedHash := CalculateNodePosition(nodeID)
			if calculatedHash != expectedHashes[i] {
				t.Errorf("Expected hash for node '%s' is '%d', but got '%d'", nodeID, expectedHashes[i], calculatedHash)
			}
		}
	})

	t.Run("KeyHashCalculation", func(t *testing.T) {
		// Test hashing for keys
		keys := []string{"key1", "key2", "key3"}
		expectedHashes := []int{
			int(crc32.ChecksumIEEE([]byte("key1"))),
			int(crc32.ChecksumIEEE([]byte("key2"))),
			int(crc32.ChecksumIEEE([]byte("key3"))),
		}

		for i, key := range keys {
			calculatedHash := CalculateKeyPosition(key)
			if calculatedHash != expectedHashes[i] {
				t.Errorf("Expected hash for key '%s' is '%d', but got '%d'", key, expectedHashes[i], calculatedHash)
			}
		}
	})

	t.Run("Consistency", func(t *testing.T) {
		// Test that the same input always produces the same hash
		nodeID := "NodeA"
		hash1 := CalculateNodePosition(nodeID)
		hash2 := CalculateNodePosition(nodeID)
		if hash1 != hash2 {
			t.Errorf("Inconsistent hash for node '%s': got '%d' and '%d'", nodeID, hash1, hash2)
		}

		key := "key1"
		hash1 = CalculateKeyPosition(key)
		hash2 = CalculateKeyPosition(key)
		if hash1 != hash2 {
			t.Errorf("Inconsistent hash for key '%s': got '%d' and '%d'", key, hash1, hash2)
		}
	})
}
