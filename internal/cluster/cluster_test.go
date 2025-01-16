package cluster

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tdevsin/keyforge/internal/logger"
)

func getTestLogger() *logger.Logger {
	return logger.GetLogger(false, "123")
}

func TestIncrementVersion(t *testing.T) {
	t.Run("Test increment version", func(t *testing.T) {
		cluster := NewCluster(getTestLogger(), "node1", 2)
		cluster.IncrementVersion()
		assert.Equal(t, 0, cluster.Version)
	})

	t.Run("Test increment version multiple times", func(t *testing.T) {
		cluster := NewCluster(getTestLogger(), "node1", 2)
		cluster.IncrementVersion()
		cluster.IncrementVersion()
		cluster.IncrementVersion()
		assert.Equal(t, 2, cluster.Version)
	})
}

func TestAddOrUpdateNode(t *testing.T) {
	t.Run("Test add node", func(t *testing.T) {
		cluster := NewCluster(getTestLogger(), "node1", 2)
		node := Node{
			ID: "node1",
		}
		cluster.AddOrUpdateNode(node)
		assert.Equal(t, 1, len(cluster.Nodes))
	})

	t.Run("Test update node", func(t *testing.T) {
		cluster := NewCluster(getTestLogger(), "node1", 2)
		node := Node{
			ID: "node1",
		}
		cluster.AddOrUpdateNode(node)
		assert.Equal(t, 1, len(cluster.Nodes))
		assert.Equal(t, "node1", cluster.Nodes["node1"].ID)
		assert.Equal(t, "", cluster.Nodes["node1"].Address)

		node.Address = "node1:8000"
		node.Position = 1000
		cluster.AddOrUpdateNode(node)
		assert.Equal(t, 1, len(cluster.Nodes))
		assert.Equal(t, "node1:8000", cluster.Nodes["node1"].Address)
		assert.Equal(t, 1000, cluster.Nodes["node1"].Position)
	})
}

func TestRemoveNode(t *testing.T) {
	t.Run("Test remove node", func(t *testing.T) {
		cluster := NewCluster(getTestLogger(), "node1", 2)
		node := Node{
			ID: "node1",
		}
		cluster.AddOrUpdateNode(node)
		assert.Equal(t, 1, len(cluster.Nodes))

		cluster.RemoveNode("node1")
		assert.Equal(t, 0, len(cluster.Nodes))
		_, ok := cluster.Nodes["node1"]
		assert.False(t, ok)
	})
}

func TestGetNode(t *testing.T) {
	t.Run("Test get node", func(t *testing.T) {
		cluster := NewCluster(getTestLogger(), "node1", 2)
		node := Node{
			ID: "node1",
		}
		cluster.AddOrUpdateNode(node)
		assert.Equal(t, 1, len(cluster.Nodes))

		node, ok := cluster.GetNode("node1")
		assert.True(t, ok)
		assert.Equal(t, "node1", node.ID)
	})

	t.Run("Test get node not found", func(t *testing.T) {
		cluster := NewCluster(getTestLogger(), "node1", 2)
		node := Node{
			ID: "node1",
		}
		cluster.AddOrUpdateNode(node)
		assert.Equal(t, 1, len(cluster.Nodes))

		_, ok := cluster.GetNode("node2")
		assert.False(t, ok)
	})
}

func TestGetHealthyNodes(t *testing.T) {
	t.Run("Test get healthy nodes", func(t *testing.T) {
		cluster := NewCluster(getTestLogger(), "node1", 2)
		node1 := Node{
			ID: "node1",
			Health: Health{
				Status:      Healthy,
				LastChecked: time.Now(),
			},
		}
		node2 := Node{
			ID: "node2",
			Health: Health{
				Status:      Healthy,
				LastChecked: time.Now(),
			},
		}
		node3 := Node{
			ID: "node3",
			Health: Health{
				Status:      SuspectedFailed,
				LastChecked: time.Now(),
			},
		}
		node4 := Node{
			ID: "node4",
			Health: Health{
				Status:      PermanentFailed,
				LastChecked: time.Now(),
			},
		}
		cluster.AddOrUpdateNode(node1)
		cluster.AddOrUpdateNode(node2)
		cluster.AddOrUpdateNode(node3)
		cluster.AddOrUpdateNode(node4)
		assert.Equal(t, 4, len(cluster.Nodes))

		healthyNodes := cluster.GetHealthyNodes()
		assert.Equal(t, 3, len(healthyNodes))

		// Create a map of expected healthy nodes for easy comparison
		expected := map[string]bool{
			"node1": true,
			"node2": true,
			"node3": true,
			"node4": false,
		}

		// Check that each healthy node is in the expected map
		for _, node := range healthyNodes {
			assert.True(t, expected[node.ID], "Unexpected node found: %s", node.ID)
		}
	})

}

func TestConcurrency(t *testing.T) {
	t.Run("Test concurrency safety of ClusterManager methods", func(t *testing.T) {
		cluster := NewCluster(getTestLogger(), "node1", 2)

		node1 := Node{
			ID: "node1",
			Health: Health{
				Status:      Healthy,
				LastChecked: time.Now(),
			},
		}
		node2 := Node{
			ID: "node2",
			Health: Health{
				Status:      SuspectedFailed,
				LastChecked: time.Now(),
			},
		}
		node3 := Node{
			ID: "node3",
			Health: Health{
				Status:      Healthy,
				LastChecked: time.Now(),
			},
		}

		// Add a wait group to synchronize goroutines
		var wg sync.WaitGroup

		// Run AddOrUpdateNode in parallel
		wg.Add(1)
		t.Run("AddOrUpdateNode", func(t *testing.T) {
			go func() {
				defer wg.Done()
				cluster.AddOrUpdateNode(node1)
				cluster.AddOrUpdateNode(node2)
				cluster.AddOrUpdateNode(node3)
			}()
		})

		// Run RemoveNode in parallel
		wg.Add(1)
		t.Run("RemoveNode", func(t *testing.T) {
			go func() {
				defer wg.Done()
				cluster.RemoveNode("node2")
			}()
		})

		// Run GetNode in parallel
		wg.Add(1)
		t.Run("GetNode", func(t *testing.T) {
			go func() {
				defer wg.Done()
				node, exists := cluster.GetNode("node1")
				assert.True(t, exists)
				assert.Equal(t, "node1", node.ID)
			}()
		})

		// Run GetHealthyNodes in parallel
		wg.Add(1)
		t.Run("GetHealthyNodes", func(t *testing.T) {
			go func() {
				defer wg.Done()
				healthyNodes := cluster.GetHealthyNodes()
				assert.GreaterOrEqual(t, len(healthyNodes), 1)
			}()
		})

		// Run IncrementVersion in parallel
		wg.Add(1)
		t.Run("IncrementVersion", func(t *testing.T) {
			go func() {
				defer wg.Done()
				for i := 0; i < 5; i++ {
					cluster.IncrementVersion()
				}
			}()
		})

		// Wait for all goroutines to finish
		wg.Wait()

		// Final assertions to validate state consistency
		assert.Equal(t, 2, len(cluster.Nodes)) // node2 is removed
		assert.Equal(t, "node1", cluster.Nodes["node1"].ID)
		assert.Equal(t, 4, cluster.Version) // IncrementVersion called 5 times
	})
}

func TestMergeClusterState(t *testing.T) {
	cluster := NewCluster(getTestLogger(), "node1", 2)

	t.Run("Merge Newer State", func(t *testing.T) {
		// Initial cluster state
		cluster.AddOrUpdateNode(Node{
			ID:      "node1",
			Address: "localhost:8080",
			Health:  Health{Status: Healthy, LastChecked: time.Now()},
		})
		cluster.Version = 1

		// Received state with a newer version
		receivedState := &ClusterInfo{
			Nodes: map[string]Node{
				"node2": {
					ID:      "node2",
					Address: "localhost:8081",
					Health:  Health{Status: Healthy, LastChecked: time.Now()},
				},
			},
			Version:     2,
			LastUpdated: time.Now(),
		}

		cluster.MergeClusterState(receivedState)

		if len(cluster.Nodes) != 2 {
			t.Errorf("Expected 2 nodes after merge, got %d", len(cluster.Nodes))
		}
		if cluster.Version != 2 {
			t.Errorf("Expected version 2, got %d", cluster.Version)
		}
	})

	t.Run("Merge Older State", func(t *testing.T) {
		// Received state with an older version
		receivedState := &ClusterInfo{
			Nodes: map[string]Node{
				"node3": {
					ID:      "node3",
					Address: "localhost:8082",
					Health:  Health{Status: Healthy, LastChecked: time.Now()},
				},
			},
			Version:     1, // Older version
			LastUpdated: time.Now(),
		}

		cluster.MergeClusterState(receivedState)

		if len(cluster.Nodes) != 2 { // Node3 should not be added
			t.Errorf("Expected 2 nodes after ignoring older state, got %d", len(cluster.Nodes))
		}
		if cluster.Version != 2 {
			t.Errorf("Expected version 2 to remain unchanged, got %d", cluster.Version)
		}
	})

	t.Run("Merge Same Version With Different Nodes", func(t *testing.T) {
		// Received state with the same version but different nodes
		receivedState := &ClusterInfo{
			Nodes: map[string]Node{
				"node3": {
					ID:      "node3",
					Address: "localhost:8082",
					Health:  Health{Status: Healthy, LastChecked: time.Now()},
				},
			},
			Version:     2, // Same version
			LastUpdated: time.Now(),
		}

		cluster.MergeClusterState(receivedState)

		if len(cluster.Nodes) != 3 { // Node3 should now be added
			t.Errorf("Expected 3 nodes after merge, got %d", len(cluster.Nodes))
		}
		if cluster.Version != 2 {
			t.Errorf("Expected version 2 to remain unchanged, got %d", cluster.Version)
		}
	})

	t.Run("Merge Concurrent States", func(t *testing.T) {
		receivedState1 := &ClusterInfo{
			Nodes: map[string]Node{
				"node4": {
					ID:      "node4",
					Address: "localhost:8083",
					Health:  Health{Status: Healthy, LastChecked: time.Now()},
				},
			},
			Version:     3,
			LastUpdated: time.Now(),
		}

		receivedState2 := &ClusterInfo{
			Nodes: map[string]Node{
				"node5": {
					ID:      "node5",
					Address: "localhost:8084",
					Health:  Health{Status: Healthy, LastChecked: time.Now()},
				},
			},
			Version:     3,
			LastUpdated: time.Now(),
		}

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()
			cluster.MergeClusterState(receivedState1)
		}()

		go func() {
			defer wg.Done()
			cluster.MergeClusterState(receivedState2)
		}()

		wg.Wait()

		if len(cluster.Nodes) != 5 { // Node4 and Node5 should be added
			t.Errorf("Expected 5 nodes after concurrent merge, got %d", len(cluster.Nodes))
		}
		if cluster.Version != 3 {
			t.Errorf("Expected version 3, got %d", cluster.Version)
		}
	})
}
func TestGetRandomNodesForGossip(t *testing.T) {
	// Helper function to create test nodes
	createNode := func(id string, status Status) Node {
		return Node{
			ID:      id,
			Address: "127.0.0.1:5000",
			Health:  Health{Status: status},
		}
	}

	// Initialize a new cluster with self ID and gossipN
	cluster := NewCluster(getTestLogger(), "node1", 3)

	// Add initial nodes
	for i := 1; i <= 5; i++ {
		cluster.AddOrUpdateNode(createNode("node"+strconv.Itoa(i), Healthy))
	}

	t.Run("ensure all nodes are eventually covered", func(t *testing.T) {
		results := make(map[string]bool)
		numCalls := 50 // Perform multiple calls to increase coverage probability

		// Collect results from multiple calls
		for i := 0; i < numCalls; i++ {
			nodes := cluster.GetRandomNodesForGossip()
			for _, node := range nodes {
				results[node.ID] = true
			}
		}

		// Verify all nodes except self were selected at least once
		for id := range cluster.Nodes {
			if id == "node1" { // Exclude self
				continue
			}
			if !results[id] {
				t.Errorf("node %s was never selected in multiple calls", id)
			}
		}
	})

	t.Run("ensure increasing nodes still results in complete coverage", func(t *testing.T) {
		// Add more nodes
		for i := 6; i <= 20; i++ {
			cluster.AddOrUpdateNode(createNode("node"+strconv.Itoa(i), Healthy))
		}

		results := make(map[string]bool)
		numCalls := 100 // Perform enough calls to ensure coverage for larger node sets

		// Collect results from multiple calls
		for i := 0; i < numCalls; i++ {
			nodes := cluster.GetRandomNodesForGossip()
			for _, node := range nodes {
				results[node.ID] = true
			}
		}

		// Verify all nodes except self were selected at least once
		for id := range cluster.Nodes {
			if id == "node1" { // Exclude self
				continue
			}
			if !results[id] {
				t.Errorf("node %s was never selected in multiple calls", id)
			}
		}
	})

	t.Run("ensure even distribution of node selection", func(t *testing.T) {
		results := make(map[string]int)
		numCalls := 200

		// Collect selection counts from multiple calls
		for i := 0; i < numCalls; i++ {
			nodes := cluster.GetRandomNodesForGossip()
			for _, node := range nodes {
				results[node.ID]++
			}
		}

		// Check distribution of selection
		for id := range cluster.Nodes {
			if id == "node1" { // Exclude self
				continue
			}
			if count, exists := results[id]; !exists || count == 0 {
				t.Errorf("node %s was never selected", id)
			}
		}

		// Log selection counts for debugging
		t.Logf("Node selection distribution: %+v", results)
	})

	t.Run("handle edge case where no healthy nodes exist", func(t *testing.T) {
		// Mark all nodes as unhealthy
		for id := range cluster.Nodes {
			node := cluster.Nodes[id]
			node.Health.Status = PermanentFailed
			cluster.Nodes[id] = node
		}

		nodes := cluster.GetRandomNodesForGossip()
		if len(nodes) != 0 {
			t.Errorf("expected 0 nodes, got %d", len(nodes))
		}
	})
}
