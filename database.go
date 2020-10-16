package whiterabbit

import (
	"strconv"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type DB struct {
	driver neo4j.Driver
}

// Open a connection to neo4j
func Open(cfg Config) (*DB, error) {

	uri := "bolt://" + cfg.GetHost() + ":" + strconv.Itoa(cfg.GetPort())
	driver, err := neo4j.NewDriver(uri,
		neo4j.BasicAuth(cfg.GetUser(), cfg.GetPassword(), ""),
		func(c *neo4j.Config) {
			c.Encrypted = cfg.GetEncrypted()
		})
	if err != nil {
		return nil, err
	}
	return &DB{driver: driver}, nil
}

// Close ...
func (db *DB) Close() error {
	return db.driver.Close()
}

// GetSession open session
func (db *DB) GetSession() (neo4j.Session, error) {
	session, err := db.driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return nil, err
	}
	return session, nil
}
