package bedrock

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	logrus_syslog "github.com/Sirupsen/logrus/hooks/syslog"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v2"
	"log/syslog"
)

// Application
type Application struct {
	*cli.App
	Options
	Config      ApplicationConfig
	ConfigBytes ConfigBytes
}

type Options struct {
	Config string
	Debug  bool
}

type ApplicationConfig struct {
	Log LogConfig `yaml:"log"`
}

type LogConfig struct {
	Level      string
	Formatter  string `yaml:"formatter"`
	SyslogName string `yaml:"syslog_name"`
}

func (lcfg *LogConfig) Load() (err error) {
	// Set log level
	var level = log.InfoLevel
	if lcfg.Level != "" {
		level, err = log.ParseLevel(lcfg.Level)
		if err != nil {
			log.Warnf("Failed to parse log level: %v", err)
		}
	}
	log.SetLevel(level)

	// Set log formatter
	var formatter = "text"
	if lcfg.Formatter != "" {
		formatter = lcfg.Formatter
	}
	logFormatter, err := ParseLogFormatter(formatter)
	if err != nil {
		return err
	}
	log.SetFormatter(logFormatter)

	// Add DebugLoggerHook if we are in debug mode
	if log.GetLevel() == log.DebugLevel {
		log.AddHook(new(DebugLoggerHook))
	}

	// Add syslog
	if lcfg.SyslogName != "" {
		syslog_hook, err := logrus_syslog.NewSyslogHook("", "", syslog.LOG_INFO, lcfg.SyslogName)
		if err == nil {
			log.AddHook(syslog_hook)
		} else {
			log.Warnf("Failed to use syslog: %v", err)
		}
	}

	return nil
}

func (app *Application) UnmarshalConfigFile(config interface{}, bytes []byte) error {
	return yaml.Unmarshal(bytes, config)
}

func (app *Application) ReadConfigFile(file string) error {
	b, err := ReadConfigFile(file)
	if err != nil {
		return err
	}

	app.ConfigBytes = b
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
	app.Options.Config = c.GlobalString("config")
	app.Options.Debug = c.GlobalBool("debug")

	return nil
}

func (app *Application) Configure() error {
	var err error

	// Set config
	if err = app.ReadConfigFile(app.Options.Config); err != nil {
		return err
	}

	// Read config
	if err = app.BindConfig(&app.Config); err != nil {
		return err
	}

	// Set log level
	var level = log.InfoLevel
	if app.Options.Debug {
		level = log.DebugLevel
	} else if app.Config.Log.Level != "" {
		level, err = log.ParseLevel(app.Config.Log.Level)
		if err != nil {
			log.Warnf("Failed to parse log level: %v", err)
		}
	}
	log.SetLevel(level)

	// Set log formatter
	var formatter = "text"
	if app.Config.Log.Formatter != "" {
		formatter = app.Config.Log.Formatter
	}
	logFormatter, err := ParseLogFormatter(formatter)
	if err != nil {
		return err
	}
	log.SetFormatter(logFormatter)

	// Add DebugLoggerHook if we are in debug mode
	if log.GetLevel() == log.DebugLevel {
		log.AddHook(new(DebugLoggerHook))
	}

	// Add syslog
	if app.Config.Log.SyslogName != "" {
		syslog_hook, err := logrus_syslog.NewSyslogHook("", "", syslog.LOG_INFO, app.Config.Log.SyslogName)
		if err == nil {
			log.AddHook(syslog_hook)
		} else {
			log.Warnf("Failed to use syslog: %v", err)
		}
	}

	return nil
}

func NewCliApp() *cli.App {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
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
	return app
}

func ParseLogFormatter(formatter string) (log.Formatter, error) {
	switch formatter {
	case "json":
		return &log.JSONFormatter{}, nil
	case "text":
		return &LogTextFormatter{}, nil
	default:
		return nil, fmt.Errorf("Unknown log formatter %s requested", formatter)
	}
}

// Creates a new application
func NewApplication() *Application {
	app := new(Application)
	app.App = NewCliApp()

	return app
}
