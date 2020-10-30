package integration

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/laurentbh/whiterabbit"
)

type user struct {
	whiterabbit.Model
	Name string
	Age  int64
}

func TestCreateFetchNode(t *testing.T) {
	LoadFixure([]string{"./fixtures/clean_all.txt"})

	neo, _ := whiterabbit.Open(Cfg{})
	defer neo.Close()
	con, _ := neo.GetConnection()
	defer con.Close()

	// create dummy user
	userName := "user " + strconv.FormatInt(rand.Int63n(100), 10)
	s := user{Name: userName}
	_, err := con.CreateNode(s)
	if err != nil {
		panic(err)
	}

	ret, err := con.FindAllNodes(user{})
	if err != nil {
		t.Errorf("findNodes %v", err)
	}
	if len(ret) != 1 {
		t.Errorf("findNodes returned too many entities")
	}
	retUser, ok := ret[0].(user)

	if ok == false {
		t.Error("findNodes return type is not a User")
	}
	if retUser.Name != userName {
		t.Error("findNodes return wrong User")
	}
}
