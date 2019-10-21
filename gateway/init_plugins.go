package gateway

import (
	"fmt"
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
	ID              string
	HomeID          string `db:"home_id" json:"homeID"`
	Name            string
	Trigger         []string
	TriggerKey      []string
	TriggerOperator []string
	TriggerValue    []string `db:"trigger_value" json:"triggerValue"`
	Action          []string
	ActionCall      []string
	ActionValue     []string `db:"action_value" json:"actionValue"`
	Status          bool
	CreatedAt       string `db:"created_at" json:"createdAt"`
	UpdatedAt       string `db:"updated_at" json:"updatedAt"`
	CreatorID       string `db:"creator_id" json:"creatorID"`
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
				// TODO: Save data get

				var device Device
				err := DB.Get(&device, `SELECT * FROM devices WHERE physical_id = $1`, id)
				if err == nil && val.FieldByName(field).String() != "" {
					row, err := DB.Exec("INSERT INTO datas (id, device_id, field, value_nbr, value_str, value_bool) VALUES ($1, $2, $3, $4, $5, $6)",
						utils.NewULID().String(),
						device.ID,
						field,
						func() float64 {
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
						func() string {
							if typeField == "string" {
								return val.FieldByName(field).String()
							}
							return ""
						}(),
						func() bool {
							if typeField == "bool" {
								return val.FieldByName(field).Interface().(bool)
							}
							return false
						}())
					fmt.Println(row)
					fmt.Println(err)
				}
			}
		}
	}

	go worker(plugin)

}

func automations() {

	// rows, err := DB.Queryx("SELECT * FROM automations WHERE UPPER(SUBSTR(trigger, INSTR(trigger, ' ')+1)) LIKE UPPER('%" + device.ID + "%')")
	// if err == nil {
	// 	fmt.Println(id)

	// 	for rows.Next() {

	// 		var auto automationStruct
	// 		err := rows.Scan(&auto.ID, &auto.HomeID, &auto.Name, pq.Array(&auto.Trigger), pq.Array(&auto.TriggerKey), pq.Array(&auto.TriggerValue), pq.Array(&auto.TriggerOperator), pq.Array(&auto.Action), pq.Array(&auto.ActionCall), pq.Array(&auto.ActionValue), &auto.Status, &auto.CreatedAt, &auto.UpdatedAt, &auto.CreatorID)
	// 		if err == nil {
	// 			count := 0

	// 			for i := 0; i < len(auto.Trigger); i++ {
	// 				var device Device
	// 				// field := val.FieldByName(utils.FindTriggerFromName(plugin.Config.Triggers, physicalName).Fields[0].Name).String()
	// 				field := val.FieldByName(auto.TriggerKey[i]).String()
	// 				// test := utils.FindTriggerFromName(plugin.Config.Triggers, physicalName).Fields[0].Name
	// 				fmt.Println("XXXX")
	// 				fmt.Println(auto.TriggerValue[i])
	// 				// fmt.Println(test)
	// 				fmt.Println(field)
	// 				fmt.Println("XXXX")

	// 				err = DB.Get(&device, `SELECT * FROM devices WHERE id = $1`, auto.Trigger[i])
	// 				if device.PhysicalID == id && auto.TriggerValue[i] == field {
	// 					count++
	// 				}
	// 			}

	// 			if count == len(auto.Trigger) {
	// 				for i := 0; i < len(auto.Action); i++ {
	// 					var device Device
	// 					err = DB.Get(&device, `SELECT * FROM devices WHERE id = $1`, auto.Action[i])
	// 					if err == nil {
	// 						pluginFromName(device.Plugin).CallAction(auto.ActionCall[i], []byte(device.Config))
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }

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
