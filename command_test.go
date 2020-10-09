package whiterabbit

import "testing"

func TestCreateNode(t *testing.T) {

	neo, _ := Open()

	type Test struct {
		Name string
		Id   int
	}
	s := Test{Name: "test attribute", Id: 415}

	err := neo.CreateNode(s)
	if err != nil {
		t.Errorf("error %s", err)
	}
	neo.Close()
}
