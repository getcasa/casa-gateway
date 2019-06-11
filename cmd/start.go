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
	OnStart    plugin.Symbol
	OnStop     plugin.Symbol
	OnData     plugin.Symbol
	CallAction plugin.Symbol
}

var plugins []Plugin

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

			onStart, err := plug.Lookup("OnStart")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			plugins[i].OnStart = onStart

			onStop, err := plug.Lookup("OnStop")
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			plugins[i].OnStop = onStop

			onData, err := plug.Lookup("OnData")
			if err == nil {
				plugins[i].OnData = onData
			}

			callAction, err := plug.Lookup("CallAction")
			if err == nil {
				plugins[i].CallAction = callAction
			}

			plugins[i].OnStart.(func())()
			fmt.Println(plugins[i].Name)
		}

		for {
			var wg sync.WaitGroup
			for _, plugin := range plugins {
				wg.Add(1)
				go func(plugin Plugin) {
					if plugin.OnData != nil {
						res := plugin.OnData.(func() interface{})()
						if res != nil {
							val := reflect.ValueOf(res).Elem()
							if val.Field(1).String() == "click" && val.Field(0).String() == "158d00019f84c6" {
								if pluginFromName("request").CallAction != nil {
									pluginFromName("request").CallAction.(func(string, []byte))("get", []byte(`{"Link": "http://192.168.1.131/toggle"}`))
								}
							} else if val.Field(1).String() == "double_click" && val.Field(0).String() == "158d00019f84c6" {
								if plugins[0].CallAction != nil {
									plugins[0].CallAction.(func(string, []byte))("get", []byte(`{"Link": "http://192.168.1.135/toggle"}`))
								}
							}
							fmt.Println(val.Field(0))
							// for i := 0; i < val.NumField(); i++ {
							// 	fmt.Println(val.Field(i))
							// }
						}
					}
					defer wg.Done()
				}(plugin)
			}
			wg.Wait()
		}
		// onStop.(func())()
	},
}
