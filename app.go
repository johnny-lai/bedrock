package bedrock

import (
	"errors"
  "fmt"
	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
	"github.com/johnny-lai/yaml"
	"io/ioutil"
	"log"
	"os"
)

type AppServicer interface {
	Configure(*Application) error
	Migrate(*Application) error
	Build(*Application) error
	Run(*Application) error
}

type Application struct {
	cli.App

	configBytes []byte

	Servicer AppServicer
	Engine   *gin.Engine

	Log func()
}

func (self *Application) ReadConfigFile(file string) error {
	if _, err := os.Stat(file); err != nil {
		return errors.New("config path not valid")
	}

	configBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	self.configBytes = configBytes
	return nil
}

func (self *Application) BindConfig(config interface{}) error {
	return yaml.Unmarshal(self.configBytes, config)
}

func (self *Application) Error(message string, err error, c *gin.Context, printStack bool, sendAirbrake bool) {
	/*
		w := gin.DefaultWriter
		w.Write([]byte(fmt.Sprintf("%s error:%v", message, err)))
		if printStack {
			trace := make([]byte, maxStackTraceSize)
			runtime.Stack(trace, false)
			w.Write([]byte(fmt.Sprintf("stack trace--\n%s\n--", trace)))
		}
		if sendAirbrake {
			airbrake.Notify(fmt.Errorf("%s error:%v", message, err), c.Request)
			defer airbrake.Flush()
		}
		c.AbortWithError(http.StatusInternalServerError, err)*/
}

func (self *Application) initCli() {
	svc := self.Servicer

	self.App = *cli.NewApp()
	self.App.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "config.yaml",
			Usage: "config file to use",
		},
	}

	self.App.Commands = []cli.Command{
		{
			Name:  "env",
			Usage: "Print the configurations",
			Action: func(c *cli.Context) {
				if err := self.ReadConfigFile(c.GlobalString("config")); err != nil {
					log.Fatal(err)
					return
				}

				var config interface{}
				if err := self.BindConfig(&config); err != nil {
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
		{
			Name:  "server",
			Usage: "Run the http server",
			Action: func(c *cli.Context) {
				if err := self.ReadConfigFile(c.GlobalString("config")); err != nil {
					log.Fatal(err)
					return
				}

				if err := svc.Configure(self); err != nil {
					log.Fatal(err)
				}

				if err := svc.Build(self); err != nil {
					log.Fatal(err)
				}

				svc.Run(self)
			},
		},
		{
			Name:  "migrate",
			Usage: "Perform database migrations",
			Action: func(c *cli.Context) {
				if err := svc.Migrate(self); err != nil {
					log.Fatal(err)
				}
			},
		},
	}
}

func NewApp(svc AppServicer) *Application {
	app := new(Application)

	app.Servicer = svc
	app.Engine = gin.Default()

	app.initCli()
	
	return app
}
