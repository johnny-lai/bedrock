package bedrock

import (
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

// ServiceApplication
type ServiceApplication struct {
	Application

	Engine *gin.Engine

	OnException func(*gin.Context, error)
}

func (app *ServiceApplication) LogException(c *gin.Context, err error) {
	log.WithFields(log.Fields{
		"trace": StackTrace(),
	}).Error(err)
}
