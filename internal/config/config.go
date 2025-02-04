package config

import (
	"os"
	"path"
	"time"

	"github.com/google/uuid"
	"github.com/tdevsin/keyforge/internal/cluster"
	"github.com/tdevsin/keyforge/internal/logger"
	"github.com/tdevsin/keyforge/internal/storage"
)

type Environment int

const (
	Dev Environment = iota
	Prod
)

type Consistency int

const (
	Strong Consistency = iota
	Eventual
)

type Config struct {
	Environment    Environment                // Environment is the environment in which the server is running
	RootDir        string                     // RootDir will contain all project related files like config, database etc.
	Logger         logger.Logging             // Logger is the instance of zap logger. This can be used for logging.
	Db             storage.Database           // Db is the instance of pebble.
	HashRing       cluster.ConsistentHashRing // HashRing stores all the nodes of the cluster in a ring
	ClusterInfo    cluster.ClusterManager     // ClusterInfo contains details of all the nodes in the cluster
	NodeInfo       *cluster.Node              // NodeInfo contains details of this node itself
	MetadataDb     storage.Database           // MetadataDb stores node related information in database for node recovery
	Consistency    Consistency                // Consistency defines if we need strong consistency or eventual consistency
	ConnectionPool *cluster.ConnectionPool    // ConnectionPool enables reusing existing connections
}

var config Config

func folderExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func ReadConfig(env Environment, nodeAddress string) *Config {
	homeDir, _ := os.UserHomeDir()
	rootDir := path.Join(homeDir, ".keyforge")
	metadataDir := path.Join(rootDir, "metadata")

	if !folderExists(rootDir) {
		os.Mkdir(rootDir, 0755)
	}

	if !folderExists(metadataDir) {
		os.Mkdir(metadataDir, 0755)
	}

	// TODO: Get real cluster information here
	id := uuid.NewString()

	metadataDb := storage.GetDatabaseInstance(logger.GetLogger(env == Prod, ""), metadataDir)
	// Check if there is any existing data about node
	v, e := metadataDb.ReadKey([]byte("node_id"))
	if e == nil {
		// Recovered Node
		id = string(v)
	} else {
		// Save current ID
		metadataDb.WriteKey([]byte("node_id"), []byte(id))
	}
	l := logger.GetLogger(env == Prod, id)
	position := cluster.CalculateNodePosition(id)
	thisNode := cluster.Node{
		ID:       id,
		Position: position,
		Address:  nodeAddress,
		Health: cluster.Health{
			Status:      cluster.Healthy,
			LastChecked: time.Now(),
		},
	}

	clusterInfo := cluster.NewCluster(l, id, 2)
	hashring := cluster.NewHashRing()
	// Allows HashRing to know when a node is added, updated or removed via the Observer interface
	clusterInfo.RegisterObserver(hashring)
	clusterInfo.AddOrUpdateNode(thisNode)
	// Register itself as an observer to listen to cluster changes and perform gossip
	clusterInfo.RegisterObserver(clusterInfo)
	// Start periodic gossip with other nodes to keep state in sync
	clusterInfo.StartPeriodicGossip()
	clusterInfo.StartPeriodicHealthCheck()

	config = Config{
		RootDir:        rootDir,
		Logger:         l,
		Db:             storage.GetDatabaseInstance(l, rootDir),
		HashRing:       hashring,
		NodeInfo:       &thisNode,
		Environment:    env,
		ClusterInfo:    clusterInfo,
		Consistency:    Strong,
		ConnectionPool: cluster.NewConnectionPool(),
	}
	return &config
}

// Cleanup closes all the resources and exits gracefully
func (c *Config) Cleanup() {
	c.Logger.Sync()
	c.Db.Close()
	c.ConnectionPool.Close()
}
