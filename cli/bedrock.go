package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/johnny-lai/bedrock"
	"github.com/johnny-lai/yaml"
	"log"
	"os"
)

var version = "unset"

func main() {
	app := cli.NewApp()
	app.Name = "bedrock"
	app.Version = version
	app.Usage = "A microservice structure for Go"
	app.Commands = []cli.Command{
		{
			Name:  "dump",
			Usage: "Reads the specified config file and prints the output",
			Action: func(c *cli.Context) {
				var config interface{}

				app := new(bedrock.Application)

				if err := app.ReadConfigFile(c.Args().First()); err != nil {
					log.Fatal(err)
					return
				}

				if err := app.BindConfig(&config); err != nil {
					log.Fatal(err)
					return
				}

				d, err := yaml.Marshal(config)
				if err != nil {
					log.Fatalf("error: %v", err)
				}
				fmt.Printf(string(d))
			},
		},
	}

	app.Run(os.Args)
}
