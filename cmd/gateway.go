package cmd

import (
	"github.com/ItsJimi/casa-gateway/gateway"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(gatewayCmd)
}

var gatewayCmd = &cobra.Command{
	Use:   "gateway",
	Short: "Start gateway to get data from external smart things",
	Long:  "Start gateway to get data from external smart things like Xiaomi Gateway, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		gateway.Start()
	},
}
