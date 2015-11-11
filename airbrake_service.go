package bedrock

import (
	"errors"
	"fmt"
	"github.com/airbrake/gobrake"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Airbrake Config
type AirbrakeConfig struct {
	ProjectID  int64
	ProjectKey string
}

// Airbrake Service. Add this to your AppServicer and call Configure and Build
// to add Airbrake support to your AppServicer.
type AirbrakeService struct {
	Config   AirbrakeConfig
	Notifier *gobrake.Notifier
}

// Configures the AirbrakeService
func (s *AirbrakeService) Configure(app *Application) error {
	return app.BindConfigAt(&s.Config, "airbrake")
}

// Prepares the application to use the AirbrakeService. The function will:
//
// 1. Add a recovery handler to gin
// 2. Replace app.OnException with a version that writes the airbrak in addition
//    to logging
// 3. Sets the Notifier object that will be used to push notices to Airbrake
func (s *AirbrakeService) Build(app *Application) error {
	s.Notifier = gobrake.NewNotifier(s.Config.ProjectID, s.Config.ProjectKey)

	app.Engine.Use(s.RecoveryMiddleware(app))
	app.OnException = s.OnException(app)
	return nil
}

// Generates a gin route handler for triggering a panic. This is used for testing
// that the recovery works.
func (s *AirbrakeService) PanicHandler(app *Application) func(*gin.Context) {
	return func(c *gin.Context) {
		panic("Panicking")
	}
}

// Generates a gin middleware for recovering from panics.
func (s *AirbrakeService) RecoveryMiddleware(app *Application) func(*gin.Context) {
	w := gin.DefaultWriter
	return func(c *gin.Context) {
		defer func() {
			if rval := recover(); rval != nil {
				rvalStr := fmt.Sprint(rval)
				w.Write([]byte(fmt.Sprintf("recovering from:%s at:%s", rvalStr, c.Request.URL)))
				err := errors.New(rvalStr)
				app.OnException(c, err)
				c.AbortWithError(http.StatusInternalServerError, err)
			}
			defer s.Notifier.Flush()
		}()
		c.Next()
	}
}

// The OnExeption replacement for app
func (s *AirbrakeService) OnException(app *Application) func(*gin.Context, error) {
	return func(c *gin.Context, err error) {
		app.LogException(c, err)
		s.Notifier.Notify(err, c.Request)
		defer s.Notifier.Flush()
	}
}
