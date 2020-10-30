package integration

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/laurentbh/whiterabbit"
)

func TestCreateFetchNode(t *testing.T) {
	LoadFixure([]string{"./fixtures/clean_all.txt"})

	neo, _ := whiterabbit.Open(Cfg{})
	defer neo.Close()
	con, _ := neo.GetConnection()
	defer con.Close()

	type User struct {
		whiterabbit.Model
		Name string
	}

	// create dummy user
	userName := "user " + strconv.FormatInt(rand.Int63n(100), 10)
	s := User{Name: userName}
	_, err := con.CreateNode(s)
	if err != nil {
		t.Errorf("error %s", err)
	}

	ret, err := con.FindNodes(User{})
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
