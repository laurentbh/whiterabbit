package whiterabbit

import (
	"fmt"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// Relation ...
type Relation struct {
	relation string
	from     interface{}
	to       interface{}
}

// MatchRelation ...
func (db *DB) MatchRelation(name string, candidates []interface{}) ([]Relation, error) {

	session, _ := db.GetSession()

	defer session.Close()

	cypher := "MATCH p=()-[r:" + name + "]->() RETURN p"
	result, err := session.Run(cypher,
		map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	if err = result.Err(); err != nil {
		return nil, err
	}
	ret := make([]Relation, 0)
	for result.Next() {
		var rel Relation
		record := result.Record()
		v := record.GetByIndex(0)
		// node := v.(neo4j.Node)
		// props := node.Props()

		// fmt.Printf("node %v\nprops %v", node, props)
		path := v.(neo4j.Path)
		for i, n := range path.Nodes() {
			fmt.Printf("\tnode %d: %#v\n", i, n)
			fmt.Printf("\t\t%#v", n.Props())
		}

		rel.from, err = ConvertNode(path.Nodes()[0], candidates)
		if err != nil {
			return ret, err
		}
		rel.to, err = ConvertNode(path.Nodes()[1], candidates)
		if err != nil {
			return ret, err
		}
		rel.relation = name
		ret = append(ret, rel)

		// nbRel := len(path.Relationships())
		// fmt.Printf("path: %#v  nb rel: %d\n", path, nbRel)
		// for _, n := range path.Relationships() {
		// 	fmt.Printf("rel : %#v\n", n)
		// }

	}
	return ret, nil

}
