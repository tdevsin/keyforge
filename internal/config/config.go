package config

import (
	"os"
	"path"

	"github.com/tdevsin/keyforge/internal/logger"
	"github.com/tdevsin/keyforge/internal/storage"
)

type Config struct {
	// RootDir will contain all project related files like config, database etc.
	RootDir string
	// Logger is the instance of zap logger. This can be used for logging.
	Logger *logger.Logger
	// Db is the instance of pebble.
	Db *storage.PebbleDB
}

var config Config

func folderExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func ReadConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	rootDir := path.Join(homeDir, ".keyforge")
	l := logger.GetLogger()
	if !folderExists(rootDir) {
		os.Mkdir(rootDir, 0755)
	}
	config = Config{
		RootDir: rootDir,
		Logger:  l,
		Db:      storage.GetDatabaseInstance(l, rootDir),
	}
	return &config
}

// Cleanup closes all the resources and exits gracefully
func (c *Config) Cleanup() {
	c.Logger.Sync()
	c.Db.Close()
}
