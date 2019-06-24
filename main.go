package main

import (
	"github.com/ItsJimi/casa-gateway/cmd"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cmd.Execute()
}
