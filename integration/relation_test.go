package integration

import (
	"testing"

	"github.com/laurentbh/whiterabbit"
	"github.com/stretchr/testify/assert"
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
	// TODO: find something to test where order is unpredictable

}
func TestRelation(t *testing.T) {
	LoadFixure([]string{
		"./fixtures/clean_all.txt",
		"./fixtures/relation_data.txt"})

	relName := "Defined_By"

	neo, _ := whiterabbit.Open(Cfg{})
	defer neo.Close()

	con, _ := neo.GetConnection()
	defer con.Close()

	relations, err := con.MatchRelation(relName, Ingredient{}, Category{})
	assert.Nil(t, err)

	for _, r := range relations {
		assert.Equal(t, relName, r.Relation)

		assert.IsType(t, Ingredient{}, r.From)
		assert.IsType(t, Category{}, r.To)
	}
}
