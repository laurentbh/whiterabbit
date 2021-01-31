package integration

import (
	"github.com/google/go-cmp/cmp"
	"math/rand"
	"strconv"
	"testing"

	"github.com/laurentbh/whiterabbit"
)

func TestConstraint(t *testing.T) {
	LoadFixure([]string{"./fixtures/clean_all.txt", "./fixtures/delete_constraint.txt"})
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
	u.Model.Labels = make(map[string]string)
	u.Model.ID = 123
	u.Model.Labels = map[string]string{"label1": "value1", "label2": "value2"}
	u.Nickname = make([]string,2)
	u.Nickname[0] = "first"
	u.Nickname[1] = "second"
	_, _, err := con.CreateNode(u)
	if err != nil {
		t.Errorf("TestCreateNode: %v", err)
	}
}
func TestDeleteNode(t *testing.T) {
	LoadFixure([]string{"./fixtures/clean_all.txt"})

	neo, _ := whiterabbit.Open(Cfg{})
	defer neo.Close()
	con, _ := neo.GetConnection()
	defer con.Close()

	u := User{Name: "user"}
	_, _, err := con.CreateNode(u)
	if err != nil {
		t.Errorf("TestDeleteNode: %v", err)
	}
	u = User{Name: "user2"}
	// whiterabbit.Model.Labels: map[string]string{"lol": "lol"}}
	u.Labels = make(map[string]string)
	u.Labels["label1"] = "value1"
	u.Labels["label2"] = "value2"

	_, _, err = con.CreateNode(u)
	if err != nil {
		t.Errorf("TestDeleteNode: %v", err)
	}
	where := map[string]interface{}{"Name": "user2"}
	ret, err := con.FindNodesClause(User{}, where, whiterabbit.StartsWith)
	if err != nil {
		t.Errorf("findNodes %v", err)
	}
	fecthed := ret[0].(User)

	err = con.DeleteNode(fecthed)
	if err != nil {
		t.Errorf("TestDeleteNode: %v", err)
	}
	ret, err = con.FindNodesClause(User{}, where, whiterabbit.StartsWith)
	if err != nil {
		t.Errorf("findNodes %v", err)
	}
	if ret != nil {
		t.Errorf("TestDeleteNode: %v", err)
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
	inUser := User{Name: userName, Age: 19, Nickname: []string{"one","two"}}
	_, _, err := con.CreateNode(inUser)
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
	if inUser.Name != retUser.Name {
		t.Errorf("TestCreateFetchNode failed on user.Name")
	}
	if inUser.Age != retUser.Age {
		t.Errorf("TestCreateFetchNode failed on user.Age")
	}
	if ! cmp.Equal(retUser.Nickname, inUser.Nickname) {
		t.Errorf("TestCreateFetchNode failed on user.Nickanme [slice of string]")
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
func TestFindNodesIgnoreCase(t *testing.T) {
	LoadFixure([]string{"./fixtures/clean_all.txt",
		"./fixtures/findNodesClause.txt"})

	neo, _ := whiterabbit.Open(Cfg{})
	defer neo.Close()
	con, _ := neo.GetConnection()
	defer con.Close()

	where := map[string]interface{}{"Name": "laurent"}
	ret, err := con.FindNodesClause(User{}, where, whiterabbit.Exact)
	if err != nil {
		t.Errorf("findNodes %v", err)
	}
	if len(ret) != 1 {
		t.Errorf("findNodes: %d elements returned, expecting %d", len(ret), 1)
	}
	for _, u := range ret {
		_, ok := u.(User)
		if ok == false {
			t.Error("findNodes return type is not a User")
		}
	}

	ret, err = con.FindNodesClause(User{}, where, whiterabbit.IgnoreCase)
	if err != nil {
		t.Errorf("findNodes %v", err)
	}
	if len(ret) != 2 {
		t.Errorf("findNodes: %d elements returned, expecting %d", len(ret), 2)
	}
	for _, u := range ret {
		_, ok := u.(User)
		if ok == false {
			t.Error("findNodes return type is not a User")
		}
	}
}
