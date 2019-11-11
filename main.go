package main

import (
	"log"

	"github.com/ItsJimi/casa-gateway/cmd"
	"github.com/ItsJimi/casa-gateway/logger"
)

func main() {
	config := logger.Configuration{
		EnableConsole:     true,
		ConsoleLevel:      logger.Debug,
		ConsoleJSONFormat: false,
		EnableFile:        true,
		FileLevel:         logger.Info,
		FileJSONFormat:    true,
		FileLocation:      "log.log",
	}

	err := logger.NewLogger(config)
	if err != nil {
		log.Fatalf("Could not instantiate log %s", err.Error())
	}

	logger.WithFields(logger.Fields{}).Debugf("Start casa")

	cmd.Execute()
}
