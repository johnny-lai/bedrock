package bedrock

import (
	"github.com/jinzhu/gorm"
)

type ConnectionHandler interface {
	DB() (*gorm.DB, error)
	Close()
}

type MySQLConnectionHandler struct {
	DSN     string
	LogMode bool
	db      *gorm.DB
}

// Returns a DB connection
func (c *MySQLConnectionHandler) DB() (*gorm.DB, error) {
	if c.db != nil {
		return c.db, nil
	}

	db, err := gorm.Open("mysql", c.DSN)
	if err != nil {
		return nil, err
	}
	db.LogMode(c.LogMode)
	c.db = db

	return db, nil
}

func (c *MySQLConnectionHandler) Close() {
	if c.db != nil {
		c.db.Close()
	}
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
