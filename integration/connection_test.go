package integration

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"

	"github.com/laurentbh/whiterabbit"
)

func TestCreateFetchNode(t *testing.T) {
	neo, _ := whiterabbit.Open(Cfg{})
	defer neo.Close()
	con, _ := neo.GetConnection()
	defer con.Close()

	type User struct {
		whiterabbit.Model
		Name string
	}

	// create dummy user
	rand := rand.Int63n(100)
	s := User{Name: "user " + strconv.FormatInt(rand, 10)}
	_, err := con.CreateNode(s)
	if err != nil {
		t.Errorf("error %s", err)
	}

	ret, err := con.FindNodes(User{})
	if err != nil {
		t.Errorf("findNodes %v", err)
	}
	fmt.Printf("--> %v", ret)
}
