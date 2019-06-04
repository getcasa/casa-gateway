package cmd

import (
	"fmt"
	"os"
	"plugin"
	"reflect"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(startCmd)
}

func println(str string) {
	fmt.Println(str)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		plug, err := plugin.Open("plugins/xiaomi/xiaomi.so")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		onStart, err := plug.Lookup("OnStart")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		onData, err := plug.Lookup("OnData")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// onStop, err := plug.Lookup("OnStop")
		// if err != nil {
		// 	fmt.Println(err)
		// 	os.Exit(1)
		// }

		onStart.(func())()
		for {
			res := onData.(func() interface{})()
			if res != nil {
				val := reflect.ValueOf(res).Elem()

				fmt.Println(val.Field(0))
			}
		}
		// onStop.(func())()
	},
}
