package whiterabbit

import (
	"strings"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

type Connection struct {
	session     neo4j.Session
	transaction neo4j.Transaction
}

func (con *Connection) SetSession(neoSession neo4j.Session) {
	con.session = neoSession
}
func (con *Connection) GetSession() neo4j.Session {
	return con.session
}

func (s *Connection) Close() {
	// TODO : return error
	s.session.Close()
	s.session = nil
}

// InTransaction ... execute given function in a transaction
func (con *Connection) InTransaction(f func(con *Connection) ([]neo4j.Result, error)) error {
	session := con.session
	trans, err := session.BeginTransaction()

	if err != nil {
		return err
	}
	con.transaction = trans
	defer func() {
		con.transaction = nil
	}()

	results, err := f(con)

	if err != nil {
		trans.Rollback()
		return err
	}
	// consume results
	for _, res := range results {
		if _, err = res.Consume(); err != nil {
			trans.Rollback()
			return err
		}
	}
	if err = trans.Commit(); err != nil {
		return err
	}
	return nil
}

// CreateNode ...
func (con *Connection) CreateNode(value interface{}) (neo4j.Result, error) {
	mapping, _ := GetMapping(value)
	cyp := createNodeCypher(mapping)

	var result neo4j.Result
	var err error
	if con.transaction == nil {
		result, err = con.session.Run(cyp, mapping.Values)
	} else {
		result, err = con.transaction.Run(cyp, mapping.Values)
	}
	if err != nil {
		return nil, err
	}
	if result.Err() != nil {
		return nil, result.Err()
	}
	return result, nil
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
