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
		var conf *config.Config
		env, _ := cmd.Flags().GetString("env")
		if env == "dev" {
			conf = config.ReadConfig(config.Dev)
		} else if env == "prod" {
			conf = config.ReadConfig(config.Prod)
		} else {
			panic("Invalid environment")
		}
		err := api.StartGRPCServer(conf)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.PersistentFlags().String("env", "dev", "Environment to run the server in")

}
