package integration

import (
	"testing"

	"github.com/neo4j/neo4j-go-driver/neo4j"

	"github.com/laurentbh/whiterabbit"
)

func TestTransactionCommit(t *testing.T) {

	neo, _ := whiterabbit.Open(whiterabbit.DefaultConfig{})
	defer neo.Close()

	con, _ := neo.GetConnection()
	defer con.Close()

	// tx := func(con *whiterabbit.Connection) ([]neo4j.Result, error) {
	// 	var results []neo4j.Result
	// 	return results, nil
	// }
	// con.InTransaction(tx)
	con.InTransaction(
		func(con *whiterabbit.Connection) ([]neo4j.Result, error) {

			cat1 := Category{Name: "cat 1"}
			cat2 := Category{Name: "cat 2"}
			cat3 := Category{Name: "cat 3"}

			var results []neo4j.Result

			r, _ := con.CreateNode(cat1)
			results = append(results, r)
			r, _ = con.CreateNode(cat2)
			results = append(results, r)
			r, _ = con.CreateNode(cat3)
			results = append(results, r)

			return results, nil
		},
	)
}
