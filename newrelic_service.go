package bedrock

import (
	"github.com/gin-gonic/contrib/newrelic"
)

// NewRelic Config
type NewRelicConfig struct {
	LicenseKey string
	AppName    string
	Verbose    bool
}

// NewRelic Service. Add this to your AppServicer and call Configure and Build
// to add NewRelic monitoring to your AppServicer
type NewRelicService struct {
	Config NewRelicConfig
}

// Configures NewRelicService
func (s *NewRelicService) Configure(app *ServiceApplication) error {
	return app.BindConfigAt(&s.Config, "newrelic")
}

// Builds NewRelicService
func (s *NewRelicService) Build(app *ServiceApplication) error {
	app.Engine.Use(newrelic.NewRelic(s.Config.LicenseKey, s.Config.AppName, s.Config.Verbose))
	return nil
}
