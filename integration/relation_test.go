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

	ret, _ := con.FindNodesClause(Ingredient{}, map[string]interface{}{"Name": "potato"}, whiterabbit.Exact)
	potato, _ := ret[0].(Ingredient)

	ret, _ = con.FindNodesClause(Ingredient{}, map[string]interface{}{"Name": "bean"}, whiterabbit.Exact)
	bean, _ := ret[0].(Ingredient)

	ret, _ = con.FindNodesClause(Category{}, map[string]interface{}{"Name": "vegetable"}, whiterabbit.Exact)
	vegetable, _ := ret[0].(Category)

	expected := []whiterabbit.Relation{
		{Relation: "Defined_By",
			From: potato,
			To:   vegetable},
		{Relation: "Like",
			From: potato,
			To:   bean},
	}

	r, err := con.RelationByNodeID(potato.ID, Ingredient{}, Category{})
	assert.Nil(t, err)
	assert.ElementsMatch(t, expected, r)
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
