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

type Config struct {
	Environment Environment                // Environment is the environment in which the server is running
	RootDir     string                     // RootDir will contain all project related files like config, database etc.
	Logger      logger.Logging             // Logger is the instance of zap logger. This can be used for logging.
	Db          storage.Database           // Db is the instance of pebble.
	HashRing    cluster.ConsistentHashRing // HashRing stores all the nodes of the cluster in a ring
	ClusterInfo cluster.ClusterManager     // ClusterInfo contains details of all the nodes in the cluster
	NodeInfo    *cluster.Node              // NodeInfo contains details of this node itself
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

	if !folderExists(rootDir) {
		os.Mkdir(rootDir, 0755)
	}

	// TODO: Get real cluster information here
	id := uuid.NewString()
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

	clusterInfo := cluster.NewCluster(id, 2)
	hashring := cluster.NewHashRing()
	// Allows HashRing to know when a node is added, updated or removed via the Observer interface
	clusterInfo.RegisterObserver(hashring)
	clusterInfo.AddOrUpdateNode(thisNode)

	l := logger.GetLogger(env == Prod, id)
	config = Config{
		RootDir:     rootDir,
		Logger:      l,
		Db:          storage.GetDatabaseInstance(l, rootDir),
		HashRing:    hashring,
		NodeInfo:    &thisNode,
		Environment: env,
		ClusterInfo: clusterInfo,
	}
	return &config
}

// Cleanup closes all the resources and exits gracefully
func (c *Config) Cleanup() {
	c.Logger.Sync()
	c.Db.Close()
}
