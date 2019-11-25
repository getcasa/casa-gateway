package cmd

import (
	"os"

	"github.com/ItsJimi/casa-gateway/gateway"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start [port]",
	Short: "Start gateway to get data from external smart things",
	Long:  "Start gateway to get data from external smart things like Xiaomi Gateway, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		port := "3000"
		if len(args) > 0 {
			port = args[0]
		}

		if os.Getenv("CASA_SERVER_PORT") == "" {
			os.Setenv("CASA_SERVER_PORT", "4353")
		}

		go gateway.StartPlugins(port)
		gateway.StartWebServer(port)
	},
}
