package integration

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/laurentbh/whiterabbit"
)

func init() {
	content, err := ioutil.ReadFile("relation_data.txt")
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

	neo, _ := whiterabbit.Open(whiterabbit.DefaultConfig{})
	defer neo.Close()

	candidate := []interface{}{Ingredient{}, Category{}}

	m, err := neo.MatchRelation("Defined_By", candidate)
	if err != nil {
		t.Errorf("call to MatchRelation: %s", err)
	}
	fmt.Printf("matches : %v", m)

}
