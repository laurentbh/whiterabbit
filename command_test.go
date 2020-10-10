package whiterabbit

import (
	"fmt"
	"testing"
)

func TestCreateNode(t *testing.T) {

	neo, _ := Open()
	defer neo.Close()

	type User struct {
		Name string
		Id   int
	}
	// s := User{Name: "user 2"}
	// err := neo.CreateNode(s)
	// if err != nil {
	// 	t.Errorf("error %s", err)
	// }

	ret, err := neo.FindNodes(User{})
	if err != nil {
		t.Errorf("findNodes %v", err)
	}
	fmt.Printf("--> %v", ret)
}
