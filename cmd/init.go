package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ItsJimi/casa-gateway/utils"
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
		data, err := json.Marshal(map[string]string{
			"id": utils.GetIDFile(),
		})
		utils.Check(err, "error")
		resp, err := http.Post("http://localhost:3000/v1/gateway", "application/json", bytes.NewReader(data))
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
		if err != nil {
			return
		}
		return
	},
}
