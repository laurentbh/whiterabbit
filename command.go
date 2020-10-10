package whiterabbit

import (
	"fmt"
	"strings"

	"github.com/laurentbh/whiterabbit/internal"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func (db *DB) FindNodes(nodeType interface{}) error {

	session, _ := db.GetSession()

	defer session.Close()
	mapping, _ := internal.GetMapping(nodeType)

	cypher := findNodeCypher(mapping)
	result, err := session.Run(cypher,
		map[string]interface{}{})
	if err != nil {
		return err
	}
	if err = result.Err(); err != nil {
		return err
	}
	for result.Next() {
		record := result.Record()
		v := record.GetByIndex(0)
		node := v.(neo4j.Node)
		props := node.Props()
		fmt.Printf("node %#v .  props %#v\n", node, props)
	}
	return nil
}
func findNodeCypher(mapping internal.Mapping) string {
	var builder strings.Builder
	builder.WriteString("MATCH (n:")
	builder.WriteString(mapping.Label)
	builder.WriteString(") RETURN n")
	return builder.String()
}

// CreateNode ...
func (db *DB) CreateNode(value interface{}) error {

	session, _ := db.GetSession()

	defer session.Close()

	mapping, _ := internal.GetMapping(value)
	cyp := createNodeCypher(mapping)

	result, err := session.Run(
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
func createNodeCypher(mapping internal.Mapping) (ret string) {

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
