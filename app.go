package bedrock

import (
  "errors"
  "fmt"
  "github.com/codegangsta/cli"
  "github.com/gin-gonic/gin"
  "gopkg.in/yaml.v1"
  "io/ioutil"
  "os"
  "log"
)

type AppServicer interface {
  Migrate(cfg map[interface{}]interface{}) error
  Build(cfg map[interface{}]interface{}, r *gin.Engine) error
}

func GetConfig(yamlPath string) (map[interface{}]interface{}, error) {
  config := make(map[interface{}]interface{})

  if _, err := os.Stat(yamlPath); err != nil {
    return config, errors.New("config path not valid")
  }

  ymlData, err := ioutil.ReadFile(yamlPath)
  if err != nil {
    return config, err
  }

  err = yaml.Unmarshal([]byte(ymlData), &config)

  for key, value := range config {
    config[key] = os.Expand(value.(string), os.Getenv)
  }
  /*config.SvcHost = os.Expand(config.SvcHost, os.Getenv)
  config.DbUser = os.Expand(config.DbUser, os.Getenv)
  config.DbPassword = os.Expand(config.DbPassword, os.Getenv)
  config.DbHost = os.Expand(config.DbHost, os.Getenv)
  config.DbName = os.Expand(config.DbName, os.Getenv)
  */

  return config, err
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
        cfg, err := GetConfig(c.GlobalString("config"))
        if err != nil {
          log.Fatal(err)
          return
        }

        d, err := yaml.Marshal(&cfg)
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
        cfg, err := GetConfig(c.GlobalString("config"))
        if err != nil {
          log.Fatal(err)
          return
        }

        r := gin.Default()
        if err = svc.Build(cfg, r); err != nil {
          log.Fatal(err)
        }

        r.Run(cfg["SvcHost"].(string))
      },
    },
    {
      Name:  "migratedb",
      Usage: "Perform database migrations",
      Action: func(c *cli.Context) {
        cfg, err := GetConfig(c.GlobalString("config"))
        if err != nil {
          log.Fatal(err)
          return
        }

        if err = svc.Migrate(cfg); err != nil {
          log.Fatal(err)
        }
      },
    },
  }

  return app
}