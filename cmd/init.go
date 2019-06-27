package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

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
		var id string

		file, err := os.OpenFile(".casa", os.O_APPEND, 0644)
		if err != nil {
			id = string(utils.NewULID().String())
			err = ioutil.WriteFile(".casa", []byte(id), 0644)
			utils.Check(err, "error")
		} else {
			data := make([]byte, 100)
			count, err := file.Read(data)
			utils.Check(err, "error")
			id = string(data[:count])
		}

		data, err := json.Marshal(map[string]string{
			"id": id,
		})
		utils.Check(err, "error")
		resp, err := http.Post("http://localhost:3000/v1/gateway", "application/json", bytes.NewReader(data))
		if err != nil {
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
