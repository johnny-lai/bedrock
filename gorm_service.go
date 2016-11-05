package bedrock

import (
	//"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
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

func (cfg *TLSConfig) Load() error {
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

		return mysql.RegisterTLSConfig(cfg.Name, &tlsConfig)
	}

	return nil
}

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
