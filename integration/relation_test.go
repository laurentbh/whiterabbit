package integration

import (
	"testing"

	"github.com/laurentbh/whiterabbit"
)

func TestRelationById(t *testing.T) {
	LoadFixure([]string{
		"./fixtures/clean_all.txt",
		"./fixtures/relation_data2.txt"})

	neo, _ := whiterabbit.Open(Cfg{})
	defer neo.Close()
	con, _ := neo.GetConnection()
	defer con.Close()

	candidate := []interface{}{Ingredient{}, Category{}}

	ret, _ := con.FindNodesClause(Ingredient{}, map[string]interface{}{"Name": "potato"}, whiterabbit.Exact)
	potato, _ := ret[0].(Ingredient)
	con.RelationByNodeID(potato.ID, candidate)
	// TODO: find somethhing to test where order is unpredictable

}
func TestRelation(t *testing.T) {
	LoadFixure([]string{
		"./fixtures/clean_all.txt",
		"./fixtures/relation_data.txt"})

	relName := "Defined_By"

	neo, _ := whiterabbit.Open(Cfg{})
	defer neo.Close()

	candidate := []interface{}{Ingredient{}, Category{}}

	con, _ := neo.GetConnection()
	defer con.Close()

	relations, err := con.MatchRelation(relName, candidate)
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
