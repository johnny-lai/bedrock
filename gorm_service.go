package bedrock

import (
	log "github.com/sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Generates a gin route handler for health checks that will ping the database
func HealthHandler(s ConnectionHandler) func(*gin.Context) {
	return func(c *gin.Context) {
		db, _ := s.DB()

		err := db.Exec("DO 1;").Error
		if err == nil {
			c.JSON(http.StatusOK, true)
		} else {
			log.WithFields(log.Fields{
				"trace": StackTrace(),
			}).Error(err)

			c.JSON(http.StatusInternalServerError, Errorf("%v", err))
		}
	}
}
