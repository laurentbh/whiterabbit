package whiterabbit

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
)

func TestCreateFetchNode(t *testing.T) {

	neo, _ := Open(DefaultConfig{})
	defer neo.Close()

	con, _ := neo.GetConnection()
	defer con.Close()

	type User struct {
		Model
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
