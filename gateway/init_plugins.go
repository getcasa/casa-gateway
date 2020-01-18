package gateway

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ItsJimi/casa-gateway/logger"
	"github.com/ItsJimi/casa-gateway/utils"
	"github.com/getcasa/sdk"
	"github.com/gorilla/websocket"
)

// LocalPlugin define structure of plugin
type LocalPlugin struct {
	Name         string
	File         string
	Config       *sdk.Configuration
	Init         func() []byte
	OnStart      func([]byte)
	OnStop       func()
	UpdateConfig func([]byte) []byte
	Discover     func() []sdk.DiscoveredDevice
	OnData       func() []sdk.Data
	CallAction   func(string, string, []byte, []byte)
	Stop         bool
}

// LocalPlugins list all loaded plugins
var LocalPlugins []LocalPlugin
var wg sync.WaitGroup
var ch chan []byte

func findPluginFile() {
	dir := "./plugins"
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".so" {
			LocalPlugins = append(LocalPlugins, LocalPlugin{
				Name: strings.Replace(info.Name(), ".so", "", -1),
				File: path,
				Stop: false,
			})
		}
		return nil
	})
	if err != nil {
		return
	}
}

func worker(localPlugin LocalPlugin) {
	if localPlugin.Stop == true {
		wg.Done()
		return
	}

	if localPlugin.OnData == nil {
		go worker(localPlugin)
		return
	}

	res := localPlugin.OnData()

	go func(res []sdk.Data) {
		var queues []Datas
		if res == nil || len(res) == 0 {
			return
		}

		for _, ressource := range res {
			for _, field := range sdk.FindDevicesFromName(localPlugin.Config.Devices, ressource.PhysicalName).Triggers {
				value := string(sdk.FindValueFromName(ressource.Values, field.Name).Value)

				if value == "" {
					continue
				}

				fmt.Println("------------")
				fmt.Println(ressource.PhysicalID)
				fmt.Println(ressource.PhysicalName)
				fmt.Println(field.Name)
				fmt.Println(field.Type)
				fmt.Println(value)
				fmt.Println("------------")

				queue := Datas{
					ID:       utils.NewULID().String(),
					DeviceID: ressource.PhysicalID,
					Field:    field.Name,
					ValueNbr: func() float64 {
						if field.Type == "int" {
							nbr, err := strconv.ParseFloat(value, 32)
							if err != nil {
								return 0
							}
							return nbr
						}
						return 0
					}(),
					ValueStr: func() string {
						if field.Type == "string" {
							return value
						}
						return ""
					}(),
					ValueBool: func() bool {
						if field.Type == "bool" {
							b, err := strconv.ParseBool(value)
							if err != nil {
								return false
							}
							return b
						}
						return false
					}(),
				}

				queues = append(queues, queue)
			}
		}

		if len(queues) == 0 {
			return
		}

		// Send Websocket message to server
		byteMessage, _ := json.Marshal(queues)
		message := WebsocketMessage{
			Action: "newData",
			Body:   byteMessage,
		}
		byteMessage, _ = json.Marshal(message)

		ch <- byteMessage

	}(res)
	go worker(localPlugin)
}

func writeMessage() {
	defer writeMessage()
	if WS == nil {
		logger.WithFields(logger.Fields{"code": "CGGIPW001"}).Errorf("%s", "Websocket is dead")
		return
	}
	err := WS.WriteMessage(websocket.TextMessage, <-ch)
	if err != nil {
		logger.WithFields(logger.Fields{"code": "CGGIPW002"}).Errorf("%s", err.Error())
	}
}

// StartPlugins load plugins
func StartPlugins(port string) {
	findPluginFile()

	wg.Add(len(LocalPlugins))

	for range time.Tick(2 * time.Second) {
		ServerIP = utils.DiscoverServer()
		if ServerIP == "" {
			logger.WithFields(logger.Fields{}).Debugf("Casa server not found, searching...")
			continue
		}
		break
	}

	for i := 0; i < len(LocalPlugins); i++ {
		go func(i int) {
			defer wg.Done()
			plug, err := plugin.Open(LocalPlugins[i].File)
			if err != nil {
				logger.WithFields(logger.Fields{"code": "CGGIPSP001"}).Errorf("%s", err.Error())
				os.Exit(1)
			}

			conf, err := plug.Lookup("Config")
			if err != nil {
				logger.WithFields(logger.Fields{"code": "CGGIPSP002"}).Errorf("%s", err.Error())
				os.Exit(1)
			}
			LocalPlugins[i].Config = conf.(*sdk.Configuration)

			init, err := plug.Lookup("Init")
			if err == nil {
				LocalPlugins[i].Init = init.(func() []byte)
			}

			updateConfig, err := plug.Lookup("UpdateConfig")
			if err == nil {
				LocalPlugins[i].UpdateConfig = updateConfig.(func([]byte) []byte)
			}

			onStart, err := plug.Lookup("OnStart")
			if err != nil {
				logger.WithFields(logger.Fields{"code": "CGGIPSP003"}).Errorf("%s", err.Error())
				os.Exit(1)
			}
			LocalPlugins[i].OnStart = onStart.(func([]byte))

			onStop, err := plug.Lookup("OnStop")
			if err != nil {
				logger.WithFields(logger.Fields{"code": "CGGIPSP004"}).Errorf("%s", err.Error())
				os.Exit(1)
			}
			LocalPlugins[i].OnStop = onStop.(func())

			onData, err := plug.Lookup("OnData")
			if err == nil {
				LocalPlugins[i].OnData = onData.(func() []sdk.Data)
			}

			discover, err := plug.Lookup("Discover")
			if err == nil {
				LocalPlugins[i].Discover = discover.(func() []sdk.DiscoveredDevice)
			}

			callAction, err := plug.Lookup("CallAction")
			if err == nil {
				LocalPlugins[i].CallAction = callAction.(func(string, string, []byte, []byte))
			}

			var statusCode int
			var plugin Plugin

			if LocalPlugins[i].Init != nil {
				// Get plugin to retrieve config from Casa server
				statusCode, plugin = GetPlugin(LocalPlugins[i].Name)

				if statusCode != 200 && statusCode != 404 {
					logger.WithFields(logger.Fields{"code": "CGGIPSP003"}).Warnf("Can't get plugin from casa server")
					return
				}
				if statusCode == 404 {
					config := LocalPlugins[i].Init()

					plugin = Plugin{
						Name:   LocalPlugins[i].Name,
						Config: string(config),
					}
					AddPlugin(plugin)
				}
			}

			LocalPlugins[i].OnStart([]byte(plugin.Config))

			logger.WithFields(logger.Fields{}).Debugf("%s launched!", LocalPlugins[i].Name)
		}(i)
	}
	wg.Wait()

	// Start websocket client
	StartWebsocketClient(port)

	ch = make(chan []byte)
	go writeMessage()

	// Start plugins workers to get data
	for _, localPlugin := range LocalPlugins {
		wg.Add(1)
		go worker(localPlugin)
	}

	wg.Wait()

	// Stop all plugins
	for _, localPlugin := range LocalPlugins {
		localPlugin.OnStop()
	}
}
