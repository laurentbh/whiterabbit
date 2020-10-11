package internal

import (
	"testing"
)

func TestConvertNode(t *testing.T) {

	mock := MockNode{
		mockLabels: []string{"TestStruct"},
		mockProps:  map[string]interface{}{"A": "123", "B": "valueForB", "C": "3.14"},
	}
	var candidate []interface{}
	candidate = append(candidate, TestStruct{})
	ret, err := convertNode(mock, candidate)
	if err != nil {
		t.Errorf("TestConvertNode : %s", err)
	}
	r, ok := ret.(TestStruct)
	if !ok {
		t.Errorf("expecting a TestStruct")
	} else {
		if r.A != 123 || r.B != "valueForB" || r.C != 3.14 {
			t.Errorf("assignement error ")
		}
	}
}

type MockNode struct {
	mockLabels []string
	mockProps  map[string]interface{}
}

func (m MockNode) Id() int64 {
	return 0
}
func (m MockNode) Labels() []string {
	return m.mockLabels
}

func (m MockNode) Props() map[string]interface{} {
	// r := make(map[string]interface{})
	// r["prop1"] = "value1"
	// r["prop2"] = "value2"
	return m.mockProps
}
