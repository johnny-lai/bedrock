package bedrock

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"gopkg.in/airbrake/gobrake.v2"
	"net/http"
)

// Airbrake Config
type AirbrakeConfig struct {
	Enabled    bool
	Host       string
	ProjectID  int64
	ProjectKey string
}

// Airbrake Service. Add this to your AppServicer and call Configure and Build
// to add Airbrake support to your AppServicer.
type AirbrakeService struct {
	Config   AirbrakeConfig
	Notifier *gobrake.Notifier
}

// Prepares the application to use the AirbrakeService. The function will:
//
// 1. Add a recovery handler to gin
// 2. Replace app.OnException with a version that writes the airbrak in addition
//    to logging
// 3. Sets the Notifier object that will be used to push notices to Airbrake
func (s *AirbrakeService) Build(r *gin.Engine) error {
	if !s.Config.Enabled {
		log.Info("AirbrakeService is disabled")
		return nil
	}

	s.Notifier = gobrake.NewNotifier(s.Config.ProjectID, s.Config.ProjectKey)
	s.Notifier.SetHost(s.Config.Host)
	s.Notifier.AddFilter(func(notice *gobrake.Notice) *gobrake.Notice {
		notice.Context["environment"] = gin.Mode()
		return notice
	})

	r.Use(s.RecoveryMiddleware())

	return nil
}

// Generates a gin route handler for triggering a panic. This is used for testing
// that the recovery works.
func (s *AirbrakeService) PanicHandler() func(*gin.Context) {
	return func(c *gin.Context) {
		panic("Panicking")
	}
}

// Generates a gin middleware for recovering from panics.
func (s *AirbrakeService) RecoveryMiddleware() func(*gin.Context) {
	return func(c *gin.Context) {
		defer func() {
			if rval := recover(); rval != nil {
				rvalStr := fmt.Sprint(rval)
				log.Errorf("recovering from:%s at:%s", rvalStr, c.Request.URL)

				err := errors.New(rvalStr)
				s.OnException(c, err)
				c.JSON(http.StatusInternalServerError, Errorf("%v", err))
			}
		}()
		c.Next()
	}
}

// The OnExeption replacement for app
func (s *AirbrakeService) OnException(c *gin.Context, err error) {
	log.WithFields(log.Fields{
		"trace": StackTrace(),
	}).Error(err)

	s.Notifier.Notify(err, c.Request)
	defer s.Notifier.Flush()
}
