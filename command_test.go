package whiterabbit

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/laurentbh/whiterabbit/internal"
	"github.com/neo4j/neo4j-go-driver/neo4j"
)

func TestCreateFetchNode(t *testing.T) {

	neo, _ := Open()
	defer neo.Close()

	type User struct {
		internal.Model
		Name string
	}

	// create dummy user
	rand := rand.Int63n(100)
	s := User{Name: "user " + strconv.FormatInt(rand, 10)}
	err := neo.CreateNode(s)
	if err != nil {
		t.Errorf("error %s", err)
	}

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
<<<<<<< HEAD
		for i, n := range path.Nodes() {
			fmt.Printf("\tnode %d: %#v\n", i, n)
			fmt.Printf("\t\t%#v", n.Props())
=======
		for _, n := range path.Nodes() {
			fmt.Printf("node : %#v\n", n)
>>>>>>> e214c170ae28cd82d2eba19330416de2a816828c
		}
		nbRel := len(path.Relationships())
		fmt.Printf("path: %#v  nb rel: %d\n", path, nbRel)
		for _, n := range path.Relationships() {
			fmt.Printf("rel : %#v\n", n)
		}

	}
}
