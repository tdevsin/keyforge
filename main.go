package main

import (
	"context"

	"github.com/tdevsin/keyforge/cmd"
	"github.com/tdevsin/keyforge/internal/config"
)

func main() {
	conf := config.ReadConfig()
	defer conf.Cleanup()
	ctx := context.WithValue(context.Background(), "config", conf)
	cmd.Execute(ctx)
}
