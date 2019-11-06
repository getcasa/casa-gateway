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

// LocalPlugin define structure of plugin
type LocalPlugin struct {
	Name       string
	File       string
	Config     *sdk.Configuration
	Init       func() []byte
	OnStart    func([]byte)
	OnStop     func()
	OnData     func() interface{}
	CallAction func(string, []byte)
	Stop       bool
}

// LocalPlugins list all loaded plugins
var LocalPlugins []LocalPlugin
var wg sync.WaitGroup

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

// PluginFromName return plugin from it name
func PluginFromName(name string) *LocalPlugin {
	for _, localPlugin := range LocalPlugins {
		if localPlugin.Name == name {
			return &localPlugin
		}
	}

	return nil
}

func worker(localPlugin LocalPlugin) {

	if localPlugin.Stop == true {
		wg.Done()
		return
	}

	if localPlugin.OnData != nil {
		res := localPlugin.OnData()

		if res != nil {
			physicalName := strings.ToLower(reflect.TypeOf(res).String()[strings.Index(reflect.TypeOf(res).String(), ".")+1:])
			val := reflect.ValueOf(res).Elem()
			id := val.FieldByName(utils.FindTriggerFromName(localPlugin.Config.Triggers, physicalName).FieldID).String()

			fmt.Println("------------")
			fmt.Println(id)
			fmt.Println(physicalName)
			for j := 0; j < val.NumField(); j++ {
				fmt.Println(val.Type().Field(j).Name)
				fmt.Println(val.Field(j))
			}
			fmt.Println("------------")

			for i := 0; i < len(utils.FindTriggerFromName(localPlugin.Config.Triggers, physicalName).Fields); i++ {

				field := utils.FindTriggerFromName(localPlugin.Config.Triggers, physicalName).Fields[i].Name
				typeField := utils.FindTriggerFromName(localPlugin.Config.Triggers, physicalName).Fields[i].Type

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
						go worker(localPlugin)
						return
					}
					err := WS.WriteMessage(websocket.TextMessage, byteMessage)
					if err != nil {
						log.Println("write:", err)
						go worker(localPlugin)
						return
					}
				}
			}
		}
	}

	go worker(localPlugin)
}

// StartPlugins load plugins
func StartPlugins(port string) {
	findPluginFile()

	start := time.Now()
	wg.Add(len(LocalPlugins))

	for {
		ServerIP = utils.DiscoverServer()
		if ServerIP == "" {
			log.Printf("Casa server not found, searching...")
			time.Sleep(time.Second * 5)
			continue
		}
		break
	}

	for i := 0; i < len(LocalPlugins); i++ {
		go func(i int) {
			defer wg.Done()
			plug, err := plugin.Open(LocalPlugins[i].File)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			conf, err := plug.Lookup("Config")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			LocalPlugins[i].Config = conf.(*sdk.Configuration)

			init, err := plug.Lookup("Init")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			LocalPlugins[i].Init = init.(func() []byte)

			onStart, err := plug.Lookup("OnStart")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			LocalPlugins[i].OnStart = onStart.(func([]byte))

			onStop, err := plug.Lookup("OnStop")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			LocalPlugins[i].OnStop = onStop.(func())

			onData, err := plug.Lookup("OnData")
			if err == nil {
				LocalPlugins[i].OnData = onData.(func() interface{})
			}

			callAction, err := plug.Lookup("CallAction")
			if err == nil {
				LocalPlugins[i].CallAction = callAction.(func(string, []byte))
			}

			// Get plugin to retrieve config from Casa server
			statusCode, plugin := GetPlugin(LocalPlugins[i].Name)

			if statusCode != 200 && statusCode != 404 {
				fmt.Println("Can't get plugin from casa server")
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

			LocalPlugins[i].OnStart([]byte(plugin.Config))

			fmt.Println(LocalPlugins[i].Name)
		}(i)
	}
	wg.Wait()
	t := time.Now()
	fmt.Println(t.Sub(start))

	// Start websocket client
	StartWebsocketClient(port)

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
