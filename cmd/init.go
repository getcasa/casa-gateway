package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ItsJimi/casa-gateway/utils"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init gateway.",
	Long:  "Init gateway.",
	Run: func(cmd *cobra.Command, args []string) {
		file, err := os.OpenFile(".casa", os.O_APPEND, 0644)
		if err != nil {
			id := []byte(utils.NewULID().String())
			errs := ioutil.WriteFile(".casa", id, 0644)
			utils.Check(errs, "error")
		} else {
			data := make([]byte, 100)
			count, err := file.Read(data)
			utils.Check(err, "error")
			fmt.Println(string(data[:count]))
		}
	},
}
