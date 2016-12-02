package bedrock

import (
	"github.com/jinzhu/gorm"
)

type ConnectionHandler interface {
	DB() (*gorm.DB, error)
	Close()
}

type SQLConnectionHandler struct {
	Dialect string
	DSN     string
	LogMode bool
	db      *gorm.DB
}

// Returns a DB connection
func (c *SQLConnectionHandler) DB() (*gorm.DB, error) {
	// TODO: Re-open DB if needed
	if c.db != nil {
		return c.db, nil
	}

	db, err := gorm.Open(c.Dialect, c.DSN)
	if err != nil {
		return nil, err
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
