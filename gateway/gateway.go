package gateway

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ItsJimi/casa-gateway/logger"
	"github.com/ItsJimi/casa-gateway/utils"
)

// ServerIP var to save IP address
var ServerIP string

// RegisterGateway register gateway in casa server
func RegisterGateway() {
	ip := utils.DiscoverServer()
	if ip == "" {
		logger.WithFields(logger.Fields{"code": "CGGGRG001"}).Errorf("Casa server not found")
		return
	}

	data, err := json.Marshal(map[string]string{
		"id": utils.GetIDFile(),
	})
	utils.Check(err, "error")
	resp, err := http.Post("http://"+ip+":"+os.Getenv("CASA_SERVER_PORT")+"/v1/gateway", "application/json", bytes.NewReader(data))
	if err != nil {
		logger.WithFields(logger.Fields{"code": "CGGGRG002"}).Errorf("%s", err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	logger.WithFields(logger.Fields{}).Debugf("body: %s", body)
	if err != nil {
		logger.WithFields(logger.Fields{"code": "CGGGRG003"}).Errorf("%s", err.Error())
		return
	}
	return
}

// GetPlugin retrieve a plugin from Casa server
func GetPlugin(pluginName string) (int, Plugin) {
	res, err := http.Get("http://" + ServerIP + ":" + os.Getenv("CASA_SERVER_PORT") + "/v1/gateway/" + utils.GetIDFile() + "/plugins/" + pluginName)
	if err != nil {
		logger.WithFields(logger.Fields{"code": "CGGGGP001"}).Errorf("%s", err.Error())
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

	res, err := http.Post("http://"+ServerIP+":"+os.Getenv("CASA_SERVER_PORT")+"/v1/gateway/"+utils.GetIDFile()+"/plugins", "application/json", bytes.NewReader(bytePlugin))
	if err != nil {
		logger.WithFields(logger.Fields{"code": "CGGGAP001"}).Errorf("%s", err.Error())
	}
	defer res.Body.Close()

	return res.StatusCode
}
