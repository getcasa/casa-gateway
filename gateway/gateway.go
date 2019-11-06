package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ItsJimi/casa-gateway/utils"
)

var ServerIP string

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

// GetPlugin retrieve a plugin from Casa server
func GetPlugin(pluginName string) (int, Plugin) {
	res, err := http.Get("http://" + ServerIP + ":" + utils.ServerPort + "/v1/gateway/" + utils.GetIDFile() + "/plugins/" + pluginName)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return res.StatusCode, Plugin{}
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, Plugin{}
	}

	var plugin Plugin
	if err := json.Unmarshal(body, &plugin); err != nil {
		return res.StatusCode, Plugin{}
	}

	return res.StatusCode, plugin
}

// AddPlugin retrieve a plugin from Casa server
func AddPlugin(plugin Plugin) int {
	bytePlugin, err := json.Marshal(plugin)
	if err != nil {
		return 0
	}

	res, err := http.Post("http://"+ServerIP+":"+utils.ServerPort+"/v1/gateway/"+utils.GetIDFile()+"/plugins", "application/json", bytes.NewReader(bytePlugin))
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	return res.StatusCode
}
