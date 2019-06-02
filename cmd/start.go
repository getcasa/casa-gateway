package cmd

import (
	"fmt"
	"os"
	"plugin"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		plug, err := plugin.Open("plugins/example/example.so")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		start, err := plug.Lookup("Start")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		start.(func())()
	},
}
