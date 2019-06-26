package cmd

import (
	"os"

	"github.com/ItsJimi/casa-gateway/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(resetCmd)
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset gateway.",
	Long:  "Reset gateway.",
	Run: func(cmd *cobra.Command, args []string) {
		err := os.Remove(".casa")
		utils.Check(err, "error")
	},
}
