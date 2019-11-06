package gateway

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"plugin"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ItsJimi/casa-gateway/utils"
	"github.com/getcasa/sdk"
	"github.com/gorilla/websocket"
)

// Plugin define structure of plugin
type Plugin struct {
	Name       string
	File       string
	Config     *sdk.Configuration
	OnStart    func()
	OnStop     func()
	OnData     func() interface{}
	CallAction func(string, []byte)
	Stop       bool
}

// Plugins list all loaded plugins
var Plugins []Plugin
var wg sync.WaitGroup

func findPluginFile() {
	dir := "./plugins"
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".so" {
			Plugins = append(Plugins, Plugin{
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

// PluginFromName return plugin from it name
func PluginFromName(name string) *Plugin {
	for _, plugin := range Plugins {
		if plugin.Name == name {
			return &plugin
		}
	}

	return nil
}

func worker(plugin Plugin) {

	if plugin.Stop == true {
		wg.Done()
		return
	}

	if plugin.OnData != nil {
		res := plugin.OnData()

		if res != nil {
			physicalName := strings.ToLower(reflect.TypeOf(res).String()[strings.Index(reflect.TypeOf(res).String(), ".")+1:])
			val := reflect.ValueOf(res).Elem()
			id := val.FieldByName(utils.FindTriggerFromName(plugin.Config.Triggers, physicalName).FieldID).String()

			fmt.Println("------------")
			fmt.Println(id)
			fmt.Println(physicalName)
			for j := 0; j < val.NumField(); j++ {
				fmt.Println(val.Type().Field(j).Name)
				fmt.Println(val.Field(j))
			}
			fmt.Println("------------")

			for i := 0; i < len(utils.FindTriggerFromName(plugin.Config.Triggers, physicalName).Fields); i++ {

				field := utils.FindTriggerFromName(plugin.Config.Triggers, physicalName).Fields[i].Name
				typeField := utils.FindTriggerFromName(plugin.Config.Triggers, physicalName).Fields[i].Type

				if val.FieldByName(field).String() != "" {
					queue := Datas{
						ID:       utils.NewULID().String(),
						DeviceID: id,
						Field:    field,
						ValueNbr: func() float64 {
							fmt.Println(val.FieldByName(field).String())
							if typeField == "int" {
								nbr, err := strconv.ParseFloat(val.FieldByName(field).String(), 32)
								if err != nil {
									return 0
								}
								return nbr
							}
							return 0
						}(),
						ValueStr: func() string {
							if typeField == "string" {
								return val.FieldByName(field).String()
							}
							return ""
						}(),
						ValueBool: func() bool {
							if typeField == "bool" {
								return val.FieldByName(field).Interface().(bool)
							}
							return false
						}(),
					}

					// Send Websocket message to server
					byteMessage, _ := json.Marshal(queue)
					message := WebsocketMessage{
						Action: "newData",
						Body:   byteMessage,
					}
					byteMessage, _ = json.Marshal(message)
					if WS == nil {
						go worker(plugin)
						return
					}
					err := WS.WriteMessage(websocket.TextMessage, byteMessage)
					if err != nil {
						log.Println("write:", err)
						go worker(plugin)
						return
					}
				}
			}
		}
	}

	go worker(plugin)
}

// StartPlugins load plugins
func StartPlugins() {
	findPluginFile()

	start := time.Now()
	wg.Add(len(Plugins))

	for i := 0; i < len(Plugins); i++ {
		go func(i int) {
			defer wg.Done()
			plug, err := plugin.Open(Plugins[i].File)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			conf, err := plug.Lookup("Config")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			Plugins[i].Config = conf.(*sdk.Configuration)

			onStart, err := plug.Lookup("OnStart")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			Plugins[i].OnStart = onStart.(func())

			onStop, err := plug.Lookup("OnStop")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			Plugins[i].OnStop = onStop.(func())

			onData, err := plug.Lookup("OnData")
			if err == nil {
				Plugins[i].OnData = onData.(func() interface{})
			}

			callAction, err := plug.Lookup("CallAction")
			if err == nil {
				Plugins[i].CallAction = callAction.(func(string, []byte))
			}

			Plugins[i].OnStart()

			fmt.Println(Plugins[i].Name)
		}(i)
	}
	wg.Wait()
	t := time.Now()
	fmt.Println(t.Sub(start))

	// Start websocket client
	StartWebsocketClient()

	// Start plugins workers to get data
	for _, plugin := range Plugins {
		wg.Add(1)
		go worker(plugin)
	}

	wg.Wait()

	// Stop all plugins
	for _, plugin := range Plugins {
		plugin.OnStop()
	}
}
