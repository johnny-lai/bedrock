package bedrock

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type DbConfig struct {
	DbUser     string
	DbPassword string
	DbHost     string
	DbName     string
}

type GormService struct {
	config DbConfig
}

func (s *GormService) Db() (gorm.DB, error) {
	connectionString := s.config.DbUser + ":" + s.config.DbPassword + "@tcp(" + s.config.DbHost + ":3306)/" + s.config.DbName + "?charset=utf8&parseTime=True"

	return gorm.Open("mysql", connectionString)
}

func (s *GormService) Configure(app *Application) error {
	return app.BindConfig(&s.config)
}

func (s *GormService) HealthHandler(app *Application) func(*gin.Context) {
	return func(c *gin.Context) {
		db, _ := s.Db()

		err := db.Exec("DO 1;")
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"status": "success"})
		} else {
			app.Error("Database execution", nil, c, true, true)
		}
	}
}
