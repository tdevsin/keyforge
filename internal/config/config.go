package config

import (
	"os"
	"path"

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
	// Environment is the environment in which the server is running
	Environment Environment
	// RootDir will contain all project related files like config, database etc.
	RootDir string
	// Logger is the instance of zap logger. This can be used for logging.
	Logger logger.Logging
	// Db is the instance of pebble.
	Db storage.Database
	// HashRing stores all the nodes of the cluster in a ring
	HashRing cluster.ConsistentHashRing
	// NodeInfo contains details of this node itself
	NodeInfo *cluster.Node
}

var config Config

func folderExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func ReadConfig(env Environment) *Config {
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
	}

	hashring := cluster.NewHashRing()
	hashring.AddNode(thisNode)
	l := logger.GetLogger(env == Prod, id)
	config = Config{
		RootDir:     rootDir,
		Logger:      l,
		Db:          storage.GetDatabaseInstance(l, rootDir),
		HashRing:    hashring,
		NodeInfo:    &thisNode,
		Environment: env,
	}
	return &config
}

// Cleanup closes all the resources and exits gracefully
func (c *Config) Cleanup() {
	c.Logger.Sync()
	c.Db.Close()
}
