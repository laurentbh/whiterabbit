package whiterabbit

import (
	"strconv"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type DB struct {
	driver neo4j.Driver
}

// Open a connection to neo4j
func Open(cfg Config) (*DB, error) {

	uri := "bolt://" + cfg.GetHost() + ":" + strconv.Itoa(cfg.GetPort())
	driver, err := neo4j.NewDriver(uri,
		neo4j.BasicAuth(cfg.GetUser(), cfg.GetPassword(), ""))
	if err != nil {
		return nil, err
	}
	return &DB{driver: driver}, nil
}

// Close ...
func (db *DB) Close() error {
	return db.driver.Close()
}

// GetConnection open and return a neo4j session
func (db *DB) GetConnection() (Connection, error) {
	session, err := db.driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return Connection{}, err
	}
	var con Connection
	con.session = session
	return con, nil
}
