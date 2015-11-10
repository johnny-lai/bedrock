package main

import (
	"github.com/codegangsta/cli"
	"github.com/johnny-lai/bedrock"
	"log"
	"os"
	"path"
	"text/template"
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
				file := c.Args().First()

				tmpl, err := template.New(path.Base(file)).ParseFiles(file)
				if err != nil {
					log.Fatal(err)
					return
				}

				tc := bedrock.TemplateContext{}
				err = tmpl.Execute(os.Stdout, &tc)
				if err != nil {
					log.Fatal(err)
				}
			},
		},
	}

	app.Run(os.Args)
}
