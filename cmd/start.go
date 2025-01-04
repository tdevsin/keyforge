package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tdevsin/keyforge/internal/api"
	"github.com/tdevsin/keyforge/internal/config"
	"github.com/tdevsin/keyforge/internal/startup"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the server at port 8080. If port is in use, it will try the next port",
	Run: func(cmd *cobra.Command, args []string) {
		var conf *config.Config
		env, _ := cmd.Flags().GetString("env")
		bootstrap, _ := cmd.Flags().GetString("bootstrap")
		address, _ := cmd.Flags().GetString("address")

		if env == "dev" {
			conf = config.ReadConfig(config.Dev, address)
		} else if env == "prod" {
			conf = config.ReadConfig(config.Prod, address)
		} else {
			panic("Invalid environment")
		}

		err := startup.StartNodeSetupInCluster(conf, bootstrap)
		if err != nil {
			panic(err)
		}

		err = api.StartGRPCServer(conf)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.PersistentFlags().StringP("env", "e", "dev", "Specifies the environment in which the server will run. Accepted values: dev, prod")
	startCmd.PersistentFlags().StringP("bootstrap", "b", "", "Specifies the address of the bootstrap node to join the cluster. Format: <host>:<port>")
	startCmd.PersistentFlags().StringP("address", "a", "", "Specifies the address of this node, used by other nodes to connect to it. This can be a DNS name or an IP address with a port. Format: <host>:<port>")

	startCmd.MarkPersistentFlagRequired("address")
}
