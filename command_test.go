package whiterabbit

import (
	"fmt"
	"testing"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func TestCreateNode(t *testing.T) {

	neo, _ := Open()
	defer neo.Close()

	type User struct {
		Name string
		Id   int
	}
	// s := User{Name: "user 2"}
	// err := neo.CreateNode(s)
	// if err != nil {
	// 	t.Errorf("error %s", err)
	// }

	ret, err := neo.FindNodes(User{})
	if err != nil {
		t.Errorf("findNodes %v", err)
	}
	fmt.Printf("--> %v", ret)
}

func TestRelation(t *testing.T) {
	neo, _ := Open()
	defer neo.Close()

	session, _ := neo.GetSession()
	defer session.Close()
	cypher := "MATCH p=()-[r:Defined_By]->() RETURN p LIMIT 25"
	result, err := session.Run(cypher,
		map[string]interface{}{})
	if err != nil {
		t.Errorf("err %s", err)
	}
	if err = result.Err(); err != nil {
		t.Errorf("err %s", err)
	}
	for result.Next() {
		record := result.Record()
		v := record.GetByIndex(0)
		// node := v.(neo4j.Node)
		path := v.(neo4j.Path)
		// props := node.Props()

		nbNodes := len(path.Nodes())
		fmt.Printf("path: %#v  nb nodes: %d\n", path, nbNodes)
		for _, n := range path.Nodes() {
			fmt.Printf("node : %#v\n", n)
		}
		nbRel := len(path.Relationships())
		fmt.Printf("path: %#v  nb rel: %d\n", path, nbRel)
		for _, n := range path.Relationships() {
			fmt.Printf("rel : %#v\n", n)
		}

	}
}
