package cmd

import (
	"net/http"

	"github.com/ItsJimi/casa-gateway/gateway"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var client http.Client

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init gateway.",
	Long:  "Init gateway.",
	Run: func(cmd *cobra.Command, args []string) {
		gateway.RegisterGateway()
	},
}
