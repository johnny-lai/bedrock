package bedrock

import (
  "errors"
  "fmt"
  "github.com/codegangsta/cli"
  "github.com/gin-gonic/gin"
  "github.com/johnny-lai/yaml"
  "io/ioutil"
  "os"
  "log"
)

type AppServicer interface {
  Config() interface{}
  Migrate() error
  Build(r *gin.Engine) error
  Run(r *gin.Engine) error
}

func ReadConfig(yamlPath string, config interface{}) error {
  if _, err := os.Stat(yamlPath); err != nil {
    return errors.New("config path not valid")
  }

  ymlData, err := ioutil.ReadFile(yamlPath)
  if err != nil {
    return err
  }

  err = yaml.Unmarshal([]byte(ymlData), config)

  return err
}

func NewApp(svc AppServicer) *cli.App {
  app := cli.NewApp()

  app.Flags = []cli.Flag{
    cli.StringFlag{
      Name:  "config, c",
      Value: "config.yaml",
      Usage: "config file to use",
    },
  }

  app.Commands = []cli.Command{
    {
      Name:  "env",
      Usage: "Print the configurations",
      Action: func(c *cli.Context) {
        err := ReadConfig(c.GlobalString("config"), svc.Config())
        if err != nil {
          log.Fatal(err)
          return
        }

        d, err := yaml.Marshal(svc.Config())
        if err != nil {
          log.Fatalf("error: %v", err)
        }
        fmt.Printf(string(d))
      },
    },
    {
      Name:  "server",
      Usage: "Run the http server",
      Action: func(c *cli.Context) {
        err := ReadConfig(c.GlobalString("config"), svc.Config())
        if err != nil {
          log.Fatal(err)
          return
        }

        r := gin.Default()
        if err = svc.Build(r); err != nil {
          log.Fatal(err)
        }

        svc.Run(r)
      },
    },
    {
      Name:  "migrate",
      Usage: "Perform database migrations",
      Action: func(c *cli.Context) {
        err := ReadConfig(c.GlobalString("config"), svc.Config())
        if err != nil {
          log.Fatal(err)
          return
        }

        if err = svc.Migrate(); err != nil {
          log.Fatal(err)
        }
      },
    },
  }

  return app
}