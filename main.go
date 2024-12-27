package main

import (
	"github.com/tdevsin/keyforge/cmd"
	"github.com/tdevsin/keyforge/internal/logger"
	"github.com/tdevsin/keyforge/internal/storage"
)

func main() {
	storage.InitializeDatabase(".")
	db := storage.GetDatabaseInstance()
	defer db.Close()
	defer logger.Sync()
	cmd.Execute()
}
