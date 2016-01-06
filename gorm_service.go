package bedrock

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

// Database Config
type DbConfig struct {
	User     string
	Password string
	Host     string
	Database string
}

// Gorm Service. Add this to your AppService and call Configure and Build to
// add Gorm DB support to your AppServicer
type GormService struct {
	Config DbConfig
}

// Returns a DB connection
func (s *GormService) Db() (gorm.DB, error) {
	connectionString := s.Config.User + ":" + s.Config.Password + "@tcp(" + s.Config.Host + ":3306)/" + s.Config.Database + "?charset=utf8&parseTime=True"

	return gorm.Open("mysql", connectionString)
}

// Configures the GormService
func (s *GormService) Configure(app *ServiceApplication) error {
	return app.BindConfigAt(&s.Config, "db")
}

// Builds GormService
func (s *GormService) Build(app *ServiceApplication) error {
	return nil
}

// Generates a gin route handler for health checks that will ping the database
func (s *GormService) HealthHandler(app *ServiceApplication) func(*gin.Context) {
	return func(c *gin.Context) {
		db, _ := s.Db()

		err := db.Exec("DO 1;").Error
		if err == nil {
			c.JSON(http.StatusOK, true)
		} else {
			app.OnException(c, err)
			c.JSON(http.StatusInternalServerError, Errorf("%v", err))
		}
	}
}
