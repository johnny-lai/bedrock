package bedrock

import (
	log "github.com/sirupsen/logrus"
	"github.com/gin-gonic/contrib/newrelic"
	"github.com/gin-gonic/gin"
)

// NewRelic Config
type NewRelicConfig struct {
	Enabled    bool
	LicenseKey string
	AppName    string
	Verbose    bool
}

// NewRelic Service. Add this to your AppServicer and call Configure and Build
// to add NewRelic monitoring to your AppServicer
type NewRelicService struct {
	Config NewRelicConfig
}

// Builds NewRelicService
func (s *NewRelicService) Build(r *gin.Engine) error {
	if !s.Config.Enabled {
		log.Info("NewRelicService is disabled")
		return nil
	}

	r.Use(newrelic.NewRelic(s.Config.LicenseKey, s.Config.AppName, s.Config.Verbose))
	return nil
}
