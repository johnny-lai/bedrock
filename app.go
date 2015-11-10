package bedrock

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"path"
	"runtime"
	"text/template"
)

// AppServicer is the expected interface of Servicer implementations.
type AppServicer interface {
	Configure(*Application) error
	Migrate(*Application) error
	Build(*Application) error
	Run(*Application) error
}

// Application
type Application struct {
	cli.App

	ConfigBytes []byte

	Servicer AppServicer
	Engine   *gin.Engine

	OnException func(*gin.Context, error)
	Log func()
}

/*
Reads the specified config file. Note that bedrock.Application will process
the config file, using text/template, with the following extra functions:

	{{.Env "ENVIRONMENT_VARIABLE"}}
	{{.Cat "File name"}}
	{{.Base64 "a string"}}
*/
func (app *Application) ReadConfigFile(file string) error {
	if _, err := os.Stat(file); err != nil {
		return errors.New("config path not valid")
	}

	tmpl, err := template.New(path.Base(file)).ParseFiles(file)
	if err != nil {
		return err
	}

	var configBytes bytes.Buffer
	tc := TemplateContext{}
	err = tmpl.Execute(&configBytes, &tc)
	if err != nil {
		return err
	}

	app.ConfigBytes = configBytes.Bytes()
	return nil
}

func (app *Application) BindConfig(config interface{}) error {
	return yaml.Unmarshal(app.ConfigBytes, config)
}

func (app *Application) BindConfigAt(config interface{}, key string) error {
	var full = make(map[interface{}]interface{})
	if err := app.BindConfig(&full); err != nil {
		log.Fatal(err)
		return err
	}
	d, err := yaml.Marshal(full[key])
	if err != nil {
		log.Fatal(err)
		return err
	}
	fmt.Printf(string(d))
	return yaml.Unmarshal([]byte(d), config)
}

func (app *Application) LogException(c *gin.Context, err error) {
	maxStackTraceSize := 4096

	w := gin.DefaultWriter
	w.Write([]byte(fmt.Sprintf("[EXCEPTION] %v\n", err)))

	trace := make([]byte, maxStackTraceSize)
	runtime.Stack(trace, false)
	w.Write([]byte(trace))
}

func (app *Application) initCli() {
	svc := app.Servicer

	app.App = *cli.NewApp()
	app.App.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "config.yaml",
			Usage: "config file to use",
		},
	}

	app.App.Commands = []cli.Command{
		{
			Name:  "env",
			Usage: "Print the configurations",
			Action: func(c *cli.Context) {
				if err := app.ReadConfigFile(c.GlobalString("config")); err != nil {
					log.Fatal(err)
					return
				}

				var config interface{}
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
		{
			Name:  "server",
			Usage: "Run the http server",
			Action: func(c *cli.Context) {
				if err := app.ReadConfigFile(c.GlobalString("config")); err != nil {
					log.Fatal(err)
					return
				}

				if err := svc.Configure(app); err != nil {
					log.Fatal(err)
				}

				if err := svc.Build(app); err != nil {
					log.Fatal(err)
				}

				svc.Run(app)
			},
		},
		{
			Name:  "migrate",
			Usage: "Perform database migrations",
			Action: func(c *cli.Context) {
				if err := svc.Migrate(app); err != nil {
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

	app.OnException = app.LogException

	app.initCli()

	return app
}
