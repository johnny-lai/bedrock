package bedrock

import (
	log "github.com/Sirupsen/logrus"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"net/url"
	"strings"
)

type ConnectionHandler interface {
	DB() (*gorm.DB, error)
	Close()
}

type SQLConnectionHandler struct {
	Dialect    string
	DSN        string
	LogMode    bool
	AutoCreate bool
	db         *gorm.DB
}

// Returns a DB connection
func (c *SQLConnectionHandler) DB() (*gorm.DB, error) {
	// TODO: Re-open DB if needed
	if c.db != nil {
		return c.db, nil
	}

	db, err := gorm.Open(c.Dialect, c.DSN)
	if err != nil {
		// MySQLError 1049 is for "unknown database"
		if myerr, ok := err.(*mysql.MySQLError); ok && myerr.Number == 1049 && c.AutoCreate {
			log.Debugf("got error %v while connecting; going to create database", err)

			if err = c.Create(); err != nil {
				return nil, err
			}

			db, err = gorm.Open(c.Dialect, c.DSN)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	db.LogMode(c.LogMode)
	c.db = db

	return db, nil
}

func (c *SQLConnectionHandler) Close() {
	if c.db != nil {
		c.db.Close()
	}
}

func (c *SQLConnectionHandler) Create() error {
	u, err := url.Parse(c.DSN)
	if err != nil {
		return err
	}

	dbName := u.Opaque[strings.Index(u.Opaque, "/")+1:]
	u.Opaque = u.Opaque[0 : strings.Index(u.Opaque, "/")+1]

	cdb, err := gorm.Open(c.Dialect, u.String())
	if err != nil {
		return err
	}
	defer cdb.Close()

	log.Infof("creating database %s", dbName)
	return cdb.Exec("CREATE DATABASE IF NOT EXISTS `" + dbName + "`").Error
}

type SimpleConnectionHandler struct {
	Database *gorm.DB
}

func (c *SimpleConnectionHandler) DB() (*gorm.DB, error) {
	return c.Database, nil
}

func (c *SimpleConnectionHandler) Close() {
	if c.Database != nil {
		c.Database.Close()
	}
}
