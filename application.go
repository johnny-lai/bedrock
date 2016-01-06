package bedrock

import (
	"bytes"
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v2"
	"os"
	"path"
	"text/template"
)

// Application
type Application struct {
	cli.App
	ConfigBytes []byte
}

// Reads the specified config file. Note that bedrock.Application will process
// the config file, using text/template, with the following extra functions:
//
//     {{.Env "ENVIRONMENT_VARIABLE"}}
//     {{.Cat "File name"}}
//     {{.Base64 "a string"}}
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

// Unmarshals the config to the specified config variable. You would use this
// to map the config to a struct. For example:
//
//     type Config struct {
//         Key1 string
//         Key2 int
//     }
//     var c Config
//     app.BindConfig(&c)
func (app *Application) BindConfig(config interface{}) error {
	return yaml.Unmarshal(app.ConfigBytes, config)
}

// Unmarshals the config at the specified key to the specified config variable.
// You would use this to map part of the config to a struct. For example:
//
//     type Config struct {
//         User string
//         Password int
//     }
//     var c Config
//     app.BindConfigAt(&c, "db")
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

	return yaml.Unmarshal([]byte(d), config)
}

func (app *Application) InitFromCliContext(c *cli.Context) error {
	// Set config
	config := c.GlobalString("config")
	if err := app.ReadConfigFile(config); err != nil {
		return err
	}

	// Set log level
	if c.GlobalBool("debug") {
		log.SetLevel(log.DebugLevel)
		log.AddHook(new(DebugLoggerHook))
	} else {
		log.SetLevel(log.InfoLevel)
	}

	return nil
}

func (app *Application) initCli() {
	app.App = *cli.NewApp()
	app.App.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "config.yaml",
			Usage: "config file to use",
		},
		cli.BoolFlag{
			Name:  "debug, d",
			Usage: "turns on debugging",
		},
	}
	app.App.Commands = []cli.Command{}
}

// Creates a new application
func NewApplication() *Application {
	app := new(Application)
	app.initCli()

	return app
}
