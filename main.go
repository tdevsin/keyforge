package main

import (
	"github.com/tdevsin/keyforge/cmd"
	"github.com/tdevsin/keyforge/internal/logger"
)

func main() {
	logger.Info("Starting KeyForge")
	cmd.Execute()
}
