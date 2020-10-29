package whiterabbit

import (
	"strings"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type Connection struct {
	session neo4j.Session
}

func (s *Connection) SetSession(neoSession neo4j.Session) {
	s.session = neoSession
}
func (s *Connection) GetSession() neo4j.Session {
	return s.session
}

func (s *Connection) Close() {
	// TODO : return error
	s.session.Close()
	s.session = nil
}

// func (db *DB) InTransactions(lambda func(db *DB)) {
// 	session, _ := db.GetConnection()

// 	lambda(db)

// 	defer session.Close()
// }
// CreateNode ...
func (con *Connection) CreateNode(value interface{}) error {
	mapping, _ := GetMapping(value)
	cyp := createNodeCypher(mapping)

	result, err := con.GetSession().Run(
		cyp,
		mapping.Values,
	)
	if err != nil {
		return err
	}
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func createNodeCypher(mapping Mapping) (ret string) {

	var builder strings.Builder
	builder.WriteString("CREATE (n:")
	builder.WriteString(mapping.Label)

	if len(mapping.Attributes) == 0 {
		builder.WriteString(")")
		return builder.String()
	}
	builder.WriteString("{")

	if len(mapping.Attributes) > 0 {
		sep := false
		for k := range mapping.Attributes {
			if sep {
				builder.WriteString(", ")
			}
			builder.WriteString(k)
			builder.WriteString(": $")
			builder.WriteString(k)
			sep = true
		}
	}
	builder.WriteString("})")
	ret = builder.String()
	return
}
func (con *Connection) FindNodes(nodeType interface{}) ([]interface{}, error) {
	mapping, _ := GetMapping(nodeType)

	cypher := findNodeCypher(mapping)
	result, err := con.GetSession().Run(cypher,
		map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if err = result.Err(); err != nil {
		return nil, err
	}
	candidate := []interface{}{nodeType}
	var ret []interface{}
	for result.Next() {
		record := result.Record()
		v := record.GetByIndex(0)
		node := v.(neo4j.Node)

		tmp, err := ConvertNode(node, candidate)
		if err != nil {
			return nil, err
		}
		ret = append(ret, tmp)
	}
	return ret, nil
}
func findNodeCypher(mapping Mapping) string {
	var builder strings.Builder
	builder.WriteString("MATCH (n:")
	builder.WriteString(mapping.Label)
	builder.WriteString(") RETURN n")
	return builder.String()
}
