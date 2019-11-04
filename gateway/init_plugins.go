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

var plugins []Plugin
var queues []Datas
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
					queue := Datas{
						ID:       utils.NewULID().String(),
						DeviceID: device.ID,
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
					if utils.FindTriggerFromName(plugin.Config.Triggers, physicalName).Fields[i].Direct {
						queues = append(queues, queue)
					}
					row, err := DB.Exec("INSERT INTO datas (id, device_id, field, value_nbr, value_str, value_bool) VALUES ($1, $2, $3, $4, $5, $6)",
						queue.ID, queue.DeviceID, queue.Field, queue.ValueNbr, queue.ValueStr, queue.ValueBool)
					fmt.Println(row)
					fmt.Println(err)
				}
			}
		}
	}

	go worker(plugin)

}

func automations() {

	for range time.Tick(250 * time.Millisecond) {
		rows, err := DB.Queryx("SELECT * FROM automations")
		if err == nil {
			for rows.Next() {
				var auto automationStruct
				err := rows.Scan(&auto.ID, &auto.HomeID, &auto.Name, pq.Array(&auto.Trigger), pq.Array(&auto.TriggerKey), pq.Array(&auto.TriggerValue), pq.Array(&auto.TriggerOperator), pq.Array(&auto.Action), pq.Array(&auto.ActionCall), pq.Array(&auto.ActionValue), &auto.Status, &auto.CreatedAt, &auto.UpdatedAt, &auto.CreatorID)
				if err == nil {
					count := 0

					for i := 0; i < len(auto.Trigger); i++ {
						var device Device
						err = DB.Get(&device, `SELECT * FROM devices WHERE id = $1`, auto.Trigger[i])
						field := utils.FindFieldFromName(utils.FindTriggerFromName(pluginFromName(device.Plugin).Config.Triggers, device.PhysicalName).Fields, auto.TriggerKey[i])

						if field.Direct {
							queue := FindDataFromID(queues, device.ID)
							if queue.DeviceID == device.ID {
								switch field.Type {
								case "string":
									if queue.ValueStr == auto.TriggerValue[i] {
										count++
									}
								case "int":
									triggerValue, err := strconv.ParseFloat(string(auto.TriggerValue[i]), 64)
									if err == nil {
										if queue.ValueNbr == triggerValue {
											count++
										}
									}
								case "bool":
								default:
								}
							}
						} else if device.ID == auto.Trigger[i] {
							var data Datas
							err = DB.Get(&data, `SELECT * FROM datas WHERE device_id = $1 AND field = $2 ORDER BY created_at DESC`, device.ID, auto.TriggerKey[i])
							switch field.Type {
							case "string":
								if data.ValueStr == auto.TriggerValue[i] {
									count++
								}
							case "int":
								firstchar := string(auto.TriggerValue[i][0])
								value, err := strconv.ParseFloat(string(auto.TriggerValue[i][1:]), 64)
								if err == nil {
									switch firstchar {
									case ">":
										if data.ValueNbr > value {
											count++
										}
									case "<":
										if data.ValueNbr < value {
											count++
										}
									case "=":
										if data.ValueNbr == value {
											count++
										}
									default:
									}
								}
							case "bool":
							default:
							}
						}
					}

					// fmt.Println(count)

					if count == len(auto.Trigger) {
						for i := 0; i < len(auto.Action); i++ {
							var device Device
							err = DB.Get(&device, `SELECT * FROM devices WHERE id = $1`, auto.Action[i])
							if err == nil {
								pluginFromName(device.Plugin).CallAction(auto.ActionCall[i], []byte(device.Config))
							}
						}
					}
				}
			}
		}
		queues = nil
	}

	go automations()
}

// FindDataFromID find data with name ID
func FindDataFromID(datas []Datas, ID string) Datas {
	for _, data := range datas {
		if data.DeviceID == ID {
			return data
		}
	}
	return Datas{}
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
	go automations()
	wg.Wait()
	for _, plugin := range plugins {
		plugin.OnStop()
	}
}
