package integration

import (
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/laurentbh/whiterabbit"
)

func loadFixure(file string) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	cypher := string(content)

	cmds := strings.Split(cypher, ";")

	neo, err := whiterabbit.Open(whiterabbit.DefaultConfig{})
	if err != nil {
		panic(err)
	}
	defer neo.Close()

	session, _ := neo.GetSession()
	defer session.Close()

	for _, c := range cmds {
		if len(c) != 0 {
			res, err := session.Run(c,
				map[string]interface{}{})
			if err != nil {
				panic(err)
			}
			if res.Err() != nil {
				panic(res.Err())

			}
		}
	}
}

func TestRelation(t *testing.T) {
	loadFixure("./relation_data.txt")

	relName := "Defined_By"

	neo, _ := whiterabbit.Open(whiterabbit.DefaultConfig{})
	defer neo.Close()

	candidate := []interface{}{Ingredient{}, Category{}}

	relations, err := neo.MatchRelation(relName, candidate)
	if err != nil {
		t.Errorf("call to MatchRelation: %s", err)
	}
	// fmt.Printf("matches : %v", relations)
	for _, r := range relations {
		if r.Relation != relName {
			t.Errorf("expected relation %s , got %s", relName, r.Relation)
		}
		if _, ok := r.From.(Ingredient); ok == false {
			t.Errorf("wrong struct in from got %v", r.From)
		}
		if _, ok := r.To.(Category); ok == false {
			t.Errorf("wrong struct in to, got %v", r.To)
		}
	}

}
