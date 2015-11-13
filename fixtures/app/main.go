package main

import (
	"github.com/johnny-lai/bedrock"
	"{{.Env "APP_NAME"}}/core/service"
	"os"
)

var version = "unset"

func main() {
	app := bedrock.NewApp(&service.Service{})
	app.Name = "{{.Env "APP_NAME"}}"
	app.Version = version
	app.Run(os.Args)
}
