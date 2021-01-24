package whiterabbit

import (
	"strconv"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// Relation ...
type Relation struct {
	Relation string
	From     interface{}
	To       interface{}
}

// RelationByNodeID return all relations for a given node id
func (con *Connection) RelationByNodeID(id int64, candidates []interface{}) ([]Relation, error) {
	cypher := "MATCH (n) –[r]-(d) WHERE id(n) = "
	cypher = cypher + strconv.FormatInt(id, 10)
	cypher = cypher + " RETURN n,r,d"

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

		rel.From, err = ConvertNode(record.GetByIndex(0).(neo4j.Node), candidates)
		if err != nil {
			return ret, err
		}
		rel.To, err = ConvertNode(record.GetByIndex(2).(neo4j.Node), candidates)
		if err != nil {
			return ret, err
		}
		neoRel, _ := record.GetByIndex(1).(neo4j.Relationship)
		rel.Relation = neoRel.Type
		ret = append(ret, rel)
	}
	return ret, nil
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

		rel.From, err = ConvertNode(path.Nodes[0], candidates)
		if err != nil {
			return ret, err
		}
		rel.To, err = ConvertNode(path.Nodes[1], candidates)
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
