package whiterabbit

import (
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// Relation ...
type Relation struct {
	Relation string
	From     interface{}
	To       interface{}
}

// MatchRelation ...
func (con *Connection) MatchRelation(name string, candidates []interface{}) ([]Relation, error) {

	cypher := "MATCH p=()-[r:" + name + "]->() RETURN p"
	result, err := con.GetSession().Run(cypher,
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

		path := v.(neo4j.Path)
		// for i, n := range path.Nodes() {
		// 	fmt.Printf("\tnode %d: %#v\n", i, n)
		// 	fmt.Printf("\t\t%#v\n", n.Props())
		// }

		rel.From, err = ConvertNode(path.Nodes()[0], candidates)
		if err != nil {
			return ret, err
		}
		rel.To, err = ConvertNode(path.Nodes()[1], candidates)
		if err != nil {
			return ret, err
		}
		rel.Relation = name
		ret = append(ret, rel)

		// nbRel := len(path.Relationships())
		// fmt.Printf("path: %#v  nb rel: %d\n", path, nbRel)
		// for _, n := range path.Relationships() {
		// 	fmt.Printf("rel : %#v\n", n)
		// }
	}
	return ret, nil
}
