package bedrock

import (
	"github.com/gin-gonic/contrib/newrelic"
)

type NewRelicConfig struct {
	LicenseKey string
	AppName    string
	Verbose    bool
}

type NewRelicService struct {
	Config NewRelicConfig
}

func (s *NewRelicService) Configure(app *Application) error {
	return app.BindConfigAt(&s.Config, "newrelic")
}

func (s *NewRelicService) Build(app *Application) error {
	app.Engine.Use(newrelic.NewRelic(s.Config.LicenseKey, s.Config.AppName, s.Config.Verbose))
	return nil
}
