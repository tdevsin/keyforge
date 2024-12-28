package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tdevsin/keyforge/internal/api"
	"github.com/tdevsin/keyforge/internal/config"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the server at port 8080. If port is in use, it will try the next port",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		conf := ctx.Value("config").(*config.Config)
		err := api.StartGRPCServer(conf)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
