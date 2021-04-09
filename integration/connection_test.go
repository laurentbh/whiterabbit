package integration

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/laurentbh/whiterabbit"
)

func TestConstraintViolation(t *testing.T) {
	LoadFixure([]string{"./fixtures/clean_all.txt", "./fixtures/constraint_violation.txt"})

	neo, _ := whiterabbit.Open(Cfg{})
	defer neo.Close()
	con, _ := neo.GetConnection()
	defer con.Close()

	type Test struct {
		Name string
	}
	_, _, err := con.CreateNode(Test{Name: "test1"})
	assert.NotNil(t, err)

	_, _, err = con.CreateNode(Test{Name: "unique"})
	assert.Nil(t, err)
}
func TestConstraintCreation(t *testing.T) {
	LoadFixure([]string{"./fixtures/clean_all.txt", "./fixtures/delete_constraint.txt"})
	neo, _ := whiterabbit.Open(Cfg{})
	defer neo.Close()
	con, _ := neo.GetConnection()
	defer con.Close()

	err := con.SetUniqueConstraint(Category{}, "Name", "cat_name_unique")
	assert.Nil(t, err, "should not return error")

	err = con.SetUniqueConstraint(Category{}, "Name", "cat_name_unique")
	assert.NotNil(t, err, "should return error")
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
	u.Nickname = make([]string, 2)
	u.Nickname[0] = "first"
	u.Nickname[1] = "second"
	_, _, err := con.CreateNode(u)

	assert.Nil(t, err)
}
func TestDeleteNode(t *testing.T) {
	LoadFixure([]string{"./fixtures/clean_all.txt"})

	neo, _ := whiterabbit.Open(Cfg{})
	defer neo.Close()
	con, _ := neo.GetConnection()
	defer con.Close()

	u := User{Name: "user"}
	_, _, err := con.CreateNode(u)
	assert.Nil(t, err)

	u = User{Name: "user2"}
	// whiterabbit.Model.Labels: map[string]string{"lol": "lol"}}
	u.Labels = make(map[string]string)
	u.Labels["label1"] = "value1"
	u.Labels["label2"] = "value2"

	_, _, err = con.CreateNode(u)
	assert.Nil(t, err)

	where := map[string]interface{}{"Name": "user2"}
	ret, err := con.FindNodesClause(User{}, where, whiterabbit.StartsWith)
	assert.Nil(t, err)

	fetched := ret[0].(User)

	err = con.DeleteNode(fetched)
	assert.Nil(t, err)

	ret, err = con.FindNodesClause(User{}, where, whiterabbit.StartsWith)

	assert.Nil(t, err)
	assert.Nil(t, ret)
}
func TestCreateFetchNode(t *testing.T) {
	LoadFixure([]string{"./fixtures/clean_all.txt"})

	neo, _ := whiterabbit.Open(Cfg{})
	defer neo.Close()
	con, _ := neo.GetConnection()
	defer con.Close()

	// create dummy user
	userName := "user " + strconv.FormatInt(rand.Int63n(100), 10)
	inUser := User{Name: userName, Age: 19, Nickname: []string{"one", "two"}}
	_, _, err := con.CreateNode(inUser)
	assert.Nil(t, err)

	ret, err := con.FindAllNodes(User{})
	assert.Nil(t, err)

	assert.Equal(t, 1, len(ret), "findNodes should return 1 element")

	assert.IsType(t, User{}, ret[0])

	assert.Equal(t, inUser.Name, (ret[0].(User)).Name)
	assert.Equal(t, inUser.Age, (ret[0].(User)).Age)
	assert.Equal(t, inUser.Nickname, (ret[0].(User)).Nickname)
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
	assert.Nil(t, err)

	assert.Equal(t, 3, len(ret))
	for _, u := range ret {
		_, ok := u.(User)
		if ok == false {
			t.Error("findNodes return type is not a User")
		}
	}
}
func TestFindByUnknownId(t *testing.T) {
	LoadFixure([]string{"./fixtures/clean_all.txt", "./fixtures/findNodeById.txt"})

	neo, _ := whiterabbit.Open(Cfg{})
	defer neo.Close()
	con, _ := neo.GetConnection()
	defer con.Close()

	type TestNode struct {
		whiterabbit.Model
		Name string
	}

	node, err := con.FindById(100000000000, TestNode{})

	assert.Nil(t, err)
	assert.Nil(t, node)
}

func TestFindById(t *testing.T) {
	LoadFixure([]string{"./fixtures/clean_all.txt", "./fixtures/findNodeById.txt"})

	neo, _ := whiterabbit.Open(Cfg{})
	defer neo.Close()
	con, _ := neo.GetConnection()
	defer con.Close()

	type TestNode struct {
		whiterabbit.Model
		Name string
	}

	where := map[string]interface{}{"Name": "node for unit test"}
	ret, err := con.FindNodesClause(TestNode{}, where, whiterabbit.Exact)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(ret))
	assert.IsType(t, TestNode{}, ret[0])

	testNodeId := ret[0].(TestNode)

	node, err := con.FindById(testNodeId.ID, TestNode{})

	assert.Nil(t, err)
	assert.IsType(t, TestNode{}, node)
	assert.Equal(t, "node for unit test", (node.(TestNode)).Name)
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
	assert.IsType(t, User{}, ret[0])
	assert.Equal(t, 2, (ret[0].(User)).Age)

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
