package integration

import (
	"errors"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"

	"github.com/laurentbh/whiterabbit"
)

func TestTransactionCommit(t *testing.T) {

	neo, _ := whiterabbit.Open(Cfg{})
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

			var results []neo4j.Result

			_, r, _ := con.CreateNode(cat1)
			results = append(results, r)
			_, r, _ = con.CreateNode(cat2)
			results = append(results, r)

			return results, nil
		},
	)
	// TODO: test the result
}

func TestTransactionRollBack(t *testing.T) {

	neo, _ := whiterabbit.Open(Cfg{})
	defer neo.Close()

	con, _ := neo.GetConnection()
	defer con.Close()

	tx := func(con *whiterabbit.Connection) ([]neo4j.Result, error) {
		cat1 := Category{Name: "cat for rollback 1"}
		cat2 := Category{Name: "cat for rollback2"}

		var results []neo4j.Result

		_, r, _ := con.CreateNode(cat1)
		results = append(results, r)
		_, r, _ = con.CreateNode(cat2)
		results = append(results, r)

		return results, errors.New("force RB")
	}
	con.InTransaction(tx)
	// TODO test results

}
