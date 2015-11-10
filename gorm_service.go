package bedrock

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type DbConfig struct {
	User     string
	Password string
	Host     string
	Database string
}

type GormService struct {
	config DbConfig
}

func (s *GormService) Db() (gorm.DB, error) {
	connectionString := s.config.User + ":" + s.config.Password + "@tcp(" + s.config.Host + ":3306)/" + s.config.Database + "?charset=utf8&parseTime=True"

	return gorm.Open("mysql", connectionString)
}

func (s *GormService) Configure(app *Application) error {
	return app.BindConfigAt(&s.config, "db")
}

func (s *GormService) HealthHandler(app *Application) func(*gin.Context) {
	return func(c *gin.Context) {
		db, _ := s.Db()

		err := db.Exec("DO 1;").Error
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"status": "success"})
		} else {
			app.OnException(c, err)
			c.AbortWithError(http.StatusInternalServerError, err)
		}
	}
}
