package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ItsJimi/casa-gateway/utils"
)

// RegisterGateway register gateway in casa server
func RegisterGateway() {
	ip := utils.DiscoverServer()
	if ip == "" {
		fmt.Println("Casa server not found")
		return
	}

	data, err := json.Marshal(map[string]string{
		"id": utils.GetIDFile(),
	})
	utils.Check(err, "error")
	resp, err := http.Post("http://"+ip+":"+utils.ServerPort+"/v1/gateway", "application/json", bytes.NewReader(data))
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
}
