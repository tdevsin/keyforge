package main

import (
	"github.com/tdevsin/keyforge/cmd"
	"github.com/tdevsin/keyforge/internal/logger"
)

func main() {
	defer logger.Sync()
	logger.Info("Starting KeyForge")
	cmd.Execute()
}
