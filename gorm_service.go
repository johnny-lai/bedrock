package bedrock

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
)

// Database Config
type DbConfig struct {
	User          string
	Password      string
	Host          string
	Port          string
	Database      string
	Encoding      string
	Tls           bool
	SslCipher     string `yaml:"sslcipher"`
	SslCaCert     string `yaml:"sslca"`
	SslClientKey  string `yaml:"sslkey"`
	SslClientCert string `yaml:"sslcert"`
}

// Gorm Service. Add this to your AppService and call Configure and Build to
// add Gorm DB support to your AppServicer
type GormService struct {
	Config DbConfig

	connectionString string
}

// Loads the config
func (s *GormService) loadConfig() error {
	var params bytes.Buffer

	cfg := s.Config

	if cfg.Port == "" {
		cfg.Port = "3306"
	}

	params.WriteString("parseTime=True&")

	if cfg.Encoding != "" {
		params.WriteString("charset=")
		params.WriteString(cfg.Encoding)
		params.WriteString("&")
	}

	if cfg.Tls {
		var tlsConfig tls.Config

		if cfg.SslCaCert != "" {
			rootCertPool := x509.NewCertPool()
			pem, err := ioutil.ReadFile(cfg.SslCaCert)
			if err != nil {
				return err
			}
			if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
				return fmt.Errorf("Failed to append PEM.")
			}

			tlsConfig.RootCAs = rootCertPool
		}

		if cfg.SslClientKey != "" && cfg.SslClientCert != "" {
			clientCert := make([]tls.Certificate, 0, 1)
			certs, err := tls.LoadX509KeyPair(cfg.SslClientKey, cfg.SslClientCert)
			if err != nil {
				return err
			}
			clientCert = append(clientCert, certs)

			tlsConfig.Certificates = clientCert
		}

		mysql.RegisterTLSConfig("custom", &tlsConfig)

		params.WriteString("tls=custom&")
	}

	s.connectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s",
		cfg.User, cfg.Password,
		cfg.Host, cfg.Port,
		cfg.Database, params)

	return nil
}

// Returns a DB connection
func (s *GormService) Db() (*gorm.DB, error) {
	return gorm.Open("mysql", s.connectionString)
}

// Configures the GormService
func (s *GormService) Configure(app *ServiceApplication) error {
	err := app.BindConfigAt(&s.Config, "db")
	if err != nil {
		return err
	}

	// Parse and fill in the defaults
	err = s.loadConfig()
	if err != nil {
		return err
	}

	return err
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
