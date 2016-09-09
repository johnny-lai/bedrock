package bedrock

import (
	//"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

type TLSConfig struct {
	Name               string
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify"`
	Cipher             string `yaml:"cipher"`
	CaCert             string `yaml:"ca"`
	ClientKey          string `yaml:"key"`
	ClientCert         string `yaml:"cert"`
}

// Gorm Service. Add this to your AppService and call Configure and Build to
// add Gorm DB support to your AppServicer
type GormService struct {
	DbConfig  mysql.Config
	TLSConfig TLSConfig

	ConnectionString string
}

func (s *GormService) loadTLSConfig(app *ServiceApplication) error {
	err := app.BindConfigAt(&s.TLSConfig, "tlsconfig")
	if err != nil {
		return err
	}

	cfg := s.TLSConfig
	if len(cfg.Name) > 0 {
		var tlsConfig tls.Config

		tlsConfig.InsecureSkipVerify = cfg.InsecureSkipVerify

		if cfg.CaCert != "" {
			rootCertPool := x509.NewCertPool()
			pem, err := ioutil.ReadFile(cfg.CaCert)
			if err != nil {
				return err
			}
			if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
				return fmt.Errorf("Failed to append PEM.")
			}

			tlsConfig.RootCAs = rootCertPool
		}

		if cfg.ClientKey != "" && cfg.ClientCert != "" {
			clientCert := make([]tls.Certificate, 0, 1)
			certs, err := tls.LoadX509KeyPair(cfg.ClientKey, cfg.ClientCert)
			if err != nil {
				return err
			}
			clientCert = append(clientCert, certs)

			tlsConfig.Certificates = clientCert
		}

		mysql.RegisterTLSConfig(cfg.Name, &tlsConfig)
	}

	return nil
}

// Loads the config
func (s *GormService) loadDbConfig(app *ServiceApplication) error {
	err := app.BindConfigAt(&s.DbConfig, "db")
	if err != nil {
		return err
	}

	s.ConnectionString = s.DbConfig.FormatDSN()

	fmt.Printf("dsnxx=%s", s.ConnectionString)
	log.Debugf("dsn=%s", s.ConnectionString)

	return nil
}

// Returns a DB connection
func (s *GormService) Db() (*gorm.DB, error) {
	return gorm.Open("mysql", s.ConnectionString)
}

// Configures the GormService
func (s *GormService) Configure(app *ServiceApplication) error {
	if err := s.loadDbConfig(app); err != nil {
		return err
	}

	if err := s.loadTLSConfig(app); err != nil {
		return err
	}

	return nil
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
