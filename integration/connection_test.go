package integration

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/laurentbh/whiterabbit"
)

func TestConstraint(t *testing.T) {
	LoadFixure([]string{"./fixtures/clean_all.txt"})
	neo, _ := whiterabbit.Open(Cfg{})
	defer neo.Close()
	con, _ := neo.GetConnection()
	defer con.Close()

	err := con.SetUniqueConstraint(Category{}, "name", "cat_name_unique")
	if err == nil {
		t.Errorf("should return error")
	}
	err = con.SetUniqueConstraint(Category{}, "Name", "cat_name_unique")
	if err != nil {
		t.Errorf("should not return error, got [%s]", err)
	}
}
func TestCreateNode(t *testing.T) {
	LoadFixure([]string{"./fixtures/clean_all.txt"})

	neo, _ := whiterabbit.Open(Cfg{})
	defer neo.Close()
	con, _ := neo.GetConnection()
	defer con.Close()

	u := User{Name: "user1"}
	_, _, err := con.CreateNode(u)
	if err != nil {
		t.Errorf("TestCreateNode: %v", err)
	}
}
func TestCreateFetchNode(t *testing.T) {
	LoadFixure([]string{"./fixtures/clean_all.txt"})

	neo, _ := whiterabbit.Open(Cfg{})
	defer neo.Close()
	con, _ := neo.GetConnection()
	defer con.Close()

	// create dummy user
	userName := "user " + strconv.FormatInt(rand.Int63n(100), 10)
	s := User{Name: userName, Age: 19}
	_, _, err := con.CreateNode(s)
	if err != nil {
		panic(err)
	}

	ret, err := con.FindAllNodes(User{})
	if err != nil {
		t.Errorf("findNodes %v", err)
	}
	if len(ret) != 1 {
		t.Errorf("findNodes returned too many entities")
	}
	retUser, ok := ret[0].(User)

	if ok == false {
		t.Error("findNodes return type is not a User")
	}
	if retUser.Name != userName {
		t.Error("findNodes return wrong User")
	}
}

func TestFindNodesClause(t *testing.T) {
	LoadFixure([]string{"./fixtures/clean_all.txt",
		"./fixtures/findNodesClause.txt"})

	neo, _ := whiterabbit.Open(Cfg{})
	defer neo.Close()
	con, _ := neo.GetConnection()
	defer con.Close()

	where := map[string]interface{}{"Name": "user"}
	ret, err := con.FindNodesClause(User{}, where, whiterabbit.StartsWith)
	if err != nil {
		t.Errorf("findNodes %v", err)
	}
	if len(ret) != 3 {
		t.Errorf("findNodes: %d elements returned, expecting %d", len(ret), 3)
	}
	for _, u := range ret {
		_, ok := u.(User)
		if ok == false {
			t.Error("findNodes return type is not a User")
		}
	}
}
func TestFindNodesMultipleClause(t *testing.T) {
	LoadFixure([]string{"./fixtures/clean_all.txt",
		"./fixtures/findNodesClause.txt"})

	neo, _ := whiterabbit.Open(Cfg{})
	defer neo.Close()
	con, _ := neo.GetConnection()
	defer con.Close()

	where := map[string]interface{}{"Name": "user", "Age": 2}
	ret, err := con.FindNodesClause(User{}, where, whiterabbit.StartsWith)
	if err != nil {
		t.Errorf("findNodes %v", err)
	}
	if len(ret) != 1 {
		t.Errorf("findNodes: %d elements returned, expecting %d", len(ret), 1)
	}
	u, ok := ret[0].(User)
	if ok == false {
		t.Error("findNodes return type is not a User")
	} else {
		if u.Age != 2 {
			t.Error("findNodes return wrong user")
		}
	}
}
