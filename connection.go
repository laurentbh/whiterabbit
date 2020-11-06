package whiterabbit

import (
	"errors"
	"strconv"
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

// Execute a cypher
func (con *Connection) Execute(cypher string, params map[string]interface{}) error {
	p := map[string]interface{}{}
	if params != nil {
		p = params
	}
	res, err := con.GetSession().Run(cypher, p)
	if err != nil {
		return err
	}
	if res.Err() != nil {
		return res.Err()
	}
	return nil
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
	mapping, err := GetMapping(value)
	if err != nil {
		return nil, err
	}
	cyp := createNodeCypher(mapping)

	var result neo4j.Result
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

// FindAllNodes finds all nodes of a given type
func (con *Connection) FindAllNodes(nodeType interface{}) ([]interface{}, error) {
	mapping, _ := GetMapping(nodeType)

	var builder strings.Builder
	builder.WriteString("MATCH (n:")
	builder.WriteString(mapping.Label)
	builder.WriteString(") RETURN n")
	return con.findNodeHelper(builder.String(), []interface{}{nodeType})
}

func (con *Connection) findNodeHelper(cypher string, candidate []interface{}) ([]interface{}, error) {
	// fmt.Println(cypher)
	result, err := con.GetSession().Run(cypher,
		map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if err = result.Err(); err != nil {
		return nil, err
	}
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

// SearchMode operation used for `WHERE` clauses
type SearchMode int

const (
	StartsWith SearchMode = 1 + iota
	Contains
	EndsWith
	Regexp
)

// FindNodesClause finds all nodes of a given type
// searchMode is applied for all string
func (con *Connection) FindNodesClause(nodeType interface{}, where map[string]interface{}, mode SearchMode) ([]interface{}, error) {
	mapping, err := GetMapping(nodeType)
	if err != nil {
		return nil, err
	}
	// TODO: must garantee that keys from where match mapping.Attributes
	var builder strings.Builder
	builder.WriteString("MATCH (n:")
	builder.WriteString(mapping.Label)
	builder.WriteString(")")
	var firstClause = true
	if len(where) != 0 {
		builder.WriteString(" WHERE ")
		for k, v := range where {
			if !firstClause {
				builder.WriteString(" AND ")
			}
			builder.WriteString("n.")
			builder.WriteString(k)
			if mapping.Attributes[k] == "string" {
				switch mode {
				case StartsWith:
					builder.WriteString(" STARTS WITH ")
				case Contains:
					builder.WriteString(" CONTAINS ")
				case EndsWith:
					builder.WriteString(" ENDS WITH ")
				case Regexp:
					builder.WriteString(" =~ ")
				}
			} else {
				builder.WriteString(" = ")
			}
			firstClause = false
			if mapping.Attributes[k] == "string" {
				builder.WriteString(" '")
				builder.WriteString(v.(string))
				builder.WriteString("'")
			} else {
				conv, err := interfaceConv(v)
				if err != nil {
					return nil, err
				}
				builder.WriteString(conv)
			}
		}
	}
	builder.WriteString(" RETURN n")

	return con.findNodeHelper(builder.String(), []interface{}{nodeType})
}
func interfaceConv(i interface{}) (string, error) {
	conv, ok := i.(int)
	if ok {
		return strconv.Itoa(conv), nil
	}
	conv64, ok := i.(int64)
	if ok {
		return strconv.FormatInt(conv64, 10), nil
	}
	return "", errors.New("interfaceConv")
}
