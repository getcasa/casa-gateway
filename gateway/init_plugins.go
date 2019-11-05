package gateway

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
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

type websocketMessage struct {
	Action string
	Body   interface{}
}

var plugins []Plugin
var wg sync.WaitGroup
var c *websocket.Conn

func findPluginFile() {
	dir := "./plugins"
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".so" {
			plugins = append(plugins, Plugin{
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

func pluginFromName(name string) *Plugin {
	for _, plugin := range plugins {
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
					message := websocketMessage{
						Action: "newData",
						Body:   queue,
					}
					byteMessage, _ := json.Marshal(message)
					err := c.WriteMessage(websocket.TextMessage, []byte(byteMessage))
					if err != nil {
						log.Println("write:", err)
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
	wg.Add(len(plugins))
	for i := 0; i < len(plugins); i++ {
		go func(i int) {
			defer wg.Done()
			plug, err := plugin.Open(plugins[i].File)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			conf, err := plug.Lookup("Config")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			plugins[i].Config = conf.(*sdk.Configuration)

			onStart, err := plug.Lookup("OnStart")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			plugins[i].OnStart = onStart.(func())

			onStop, err := plug.Lookup("OnStop")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			plugins[i].OnStop = onStop.(func())

			onData, err := plug.Lookup("OnData")
			if err == nil {
				plugins[i].OnData = onData.(func() interface{})
			}

			callAction, err := plug.Lookup("CallAction")
			if err == nil {
				plugins[i].CallAction = callAction.(func(string, []byte))
			}

			plugins[i].OnStart()
			fmt.Println(plugins[i].Name)
		}(i)
	}
	wg.Wait()
	t := time.Now()
	fmt.Println(t.Sub(start))

	u := url.URL{Scheme: "ws", Host: "192.168.1.21:3000", Path: "/v1/ws"}
	log.Printf("connecting to %s", u.String())
	var err error
	c, _, err = websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	err = c.WriteMessage(websocket.TextMessage, []byte("Hello from Gateway"))
	if err != nil {
		log.Println("write:", err)
		return
	}

	for _, plugin := range plugins {
		wg.Add(1)
		go worker(plugin)
	}
	wg.Wait()
	for _, plugin := range plugins {
		plugin.OnStop()
	}
}
