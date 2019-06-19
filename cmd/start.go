package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"reflect"
	"strings"
	"sync"

	"github.com/getcasa/sdk"
	"github.com/spf13/cobra"
)

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

func init() {
	rootCmd.AddCommand(startCmd)
}

func println(str string) {
	fmt.Println(str)
}

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
			val := reflect.ValueOf(res).Elem()
			if val.Field(1).String() == "click" && val.Field(0).String() == "158d00019f84c6" {
				if pluginFromName("request").CallAction != nil {
					pluginFromName("request").CallAction("get", []byte(`{"Link": "http://192.168.1.131/toggle"}`))
				}
			} else if val.Field(1).String() == "double_click" && val.Field(0).String() == "158d00019f84c6" {
				if plugins[0].CallAction != nil {
					plugins[0].CallAction("get", []byte(`{"Link": "http://192.168.1.135/toggle"}`))
				}
			}
			for i := 0; i < val.NumField(); i++ {
				fmt.Println(val.Type().Field(i).Name)
				fmt.Println(val.Field(i))
			}
		}
	}
	go worker(plugin)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start gateway to get data from external smart things",
	Long:  "Start gateway to get data from external smart things like Xiaomi Gateway, etc.",
	Run: func(cmd *cobra.Command, args []string) {
		findPluginFile()

		for i := 0; i < len(plugins); i++ {
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

			// swi, err := plug.Lookup("Test")
			// if err != nil {
			// 	fmt.Println(err)
			// 	os.Exit(1)
			// }

			// fmt.Println(swi)

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
		}

		for _, plugin := range plugins {
			wg.Add(1)
			go worker(plugin)
		}
		wg.Wait()
		for _, plugin := range plugins {
			plugin.OnStop()
		}
	},
}
