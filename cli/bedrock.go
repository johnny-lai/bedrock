package main

import (
  "bufio"
  "fmt"
  "github.com/johnny-lai/bedrock"
  "github.com/johnny-lai/yaml"
  "github.com/codegangsta/cli"
  "os"
  "log"
)

var version = "unset"

func ExpandFile(file string) {
}

func main() {
  app := cli.NewApp()
  app.Name = "bedrock"
  app.Version = version
  app.Usage = "A microservice structure for Go"
  app.Commands = []cli.Command{
    {
      Name: "dump",
      Usage: "Reads the specified config file and prints the output",
      Action: func(c *cli.Context) {
        var config interface{}

        err := bedrock.ReadConfig(c.Args().First(), &config)
        if err != nil {
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
