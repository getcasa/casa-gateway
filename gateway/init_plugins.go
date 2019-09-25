package gateway

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/ItsJimi/casa-gateway/utils"
	"github.com/getcasa/sdk"
	"github.com/lib/pq"
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

var plugins []Plugin
var wg sync.WaitGroup

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

type automationStruct struct {
	ID           string
	HomeID       string `db:"home_id" json:"homeID"`
	Name         string
	Trigger      []string
	TriggerValue []string `db:"trigger_value" json:"triggerValue"`
	Action       []string
	ActionValue  []string `db:"action_value" json:"actionValue"`
	Status       bool
	CreatedAt    string `db:"created_at" json:"createdAt"`
	CreatorID    string `db:"creator_id" json:"creatorID"`
}

func worker(plugin Plugin) {

	if plugin.Stop == true {
		wg.Done()
		return
	}

	if plugin.OnData != nil {
		res := plugin.OnData()

		if res != nil {
			fmt.Println(plugin.Name)
			physicalName := strings.ToLower(reflect.TypeOf(res).String()[strings.Index(reflect.TypeOf(res).String(), ".")+1:])
			val := reflect.ValueOf(res).Elem()
			id := val.FieldByName(utils.FindTriggerFromName(plugin.Config.Triggers, physicalName).FieldID).String()
			field := val.FieldByName(utils.FindTriggerFromName(plugin.Config.Triggers, physicalName).Field).String()
			// TODO: Save data get

			fmt.Println("------------")
			for i := 0; i < val.NumField(); i++ {
				fmt.Println(val.Type().Field(i).Name)
				fmt.Println(val.Field(i))
			}
			fmt.Println("------------")

			rows, err := DB.Queryx("SELECT * FROM automations WHERE UPPER(SUBSTR(trigger, INSTR(trigger, ' ')+1)) LIKE UPPER('%" + id + "%')")
			fmt.Println(id)
			if err == nil {

				for rows.Next() {

					var auto automationStruct
					err := rows.Scan(&auto.ID, &auto.HomeID, &auto.Name, pq.Array(&auto.Trigger), pq.Array(&auto.TriggerValue), pq.Array(&auto.Action), pq.Array(&auto.ActionValue), &auto.Status, &auto.CreatedAt, &auto.CreatorID)
					if err == nil {
						count := 0

						for i := 0; i < len(auto.Trigger); i++ {
							if auto.Trigger[i] == id && auto.TriggerValue[i] == field {
								count++
							}
						}

						if count == len(auto.Trigger) {
							for i := 0; i < len(auto.Action); i++ {
								var device Device
								err = DB.Get(&device, `SELECT * FROM devices WHERE physical_id = $1`, auto.Action[i])
								if err == nil {
									pluginFromName(device.Plugin).CallAction(device.PhysicalName, []byte(`{"address": "`+auto.Action[i]+`"}`))
								}
							}
						}
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

	for _, plugin := range plugins {
		wg.Add(1)
		go worker(plugin)
	}
	wg.Wait()
	for _, plugin := range plugins {
		plugin.OnStop()
	}
}
