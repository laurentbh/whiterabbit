package whiterabbit

import (
	"strconv"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type DB struct {
	driver neo4j.Driver
}

const (
	host      = "localhost"
	port      = 7687
	encrypted = false
	user      = "root"
	password  = "neo4j"
)

// Open a connection to neo4j
// TODO : pass a config
func Open() (*DB, error) {

	uri := "bolt://" + host + ":" + strconv.Itoa(port)
	driver, err := neo4j.NewDriver(uri,
		neo4j.BasicAuth(user, password, ""),
		func(c *neo4j.Config) {
			c.Encrypted = encrypted
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
func (r *DB) GetSession() (neo4j.Session, error) {
	session, err := r.driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return nil, err
	}
	return session, nil
}
