package bedrock

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

// ServiceApplication
type ServiceApplication struct {
	Application

	Servicer AppServicer
	Engine   *gin.Engine

	OnException func(*gin.Context, error)
}

// AppServicer is the expected interface of Servicer implementations.
type AppServicer interface {
	Configure(*ServiceApplication) error
	Migrate(*ServiceApplication) error
	Build(*ServiceApplication) error
	Run(*ServiceApplication) error
}

func (app *ServiceApplication) initCli() {
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
				if err := app.InitFromCliContext(c); err != nil {
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
				if err := app.InitFromCliContext(c); err != nil {
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

func (app *ServiceApplication) LogException(c *gin.Context, err error) {
	log.WithFields(log.Fields{
		"trace": StackTrace(),
	}).Error(err)
}

// Creates a new application with the specified AppServicer
func NewServiceApplication(svc AppServicer) *ServiceApplication {
	app := new(ServiceApplication)

	app.Servicer = svc
	app.Engine = gin.Default()

	app.OnException = app.LogException

	app.initCli()

	return app
}
