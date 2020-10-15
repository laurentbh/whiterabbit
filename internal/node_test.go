package internal

import (
	"testing"
)

func TestConvertNode(t *testing.T) {

	mock := MockNode{
		id:         6,
		mockLabels: []string{"TestStruct"},
		mockProps: map[string]interface{}{
			"A":        "123",
			"B":        "valueForB",
			"C":        "3.14",
			"newLabel": "newValue",
			"label2":   "label2Value"},
	}
	var candidate []interface{}
	candidate = append(candidate, TestStruct{})
	ret, err := ConvertNode(mock, candidate)
	if err != nil {
		t.Errorf("TestConvertNode : %s", err)
	}
	r, ok := ret.(TestStruct)
	if !ok {
		t.Errorf("expecting a TestStruct")
	} else {
		if r.A != 123 || r.B != "valueForB" || r.C != 3.14 {
			t.Errorf("assignement error ")
		} else {
			// testing model
			if r.Model.ID != 6 {
				t.Errorf("model id error")
			}
			tmpV, ok := r.Labels["newLabel"]
			if ok == false {
				t.Errorf("newLabel key missing")
			}
			if tmpV != "newValue" {
				t.Errorf("model.Labels[\"newLabel\"] is %s, should be newValue", tmpV)
			}
			tmpV, ok = r.Labels["label2"]
			if ok == false {
				t.Errorf("label2 key missing")
			}
			if tmpV != "label2Value" {
				t.Errorf("model.Labels[\"label2\"] is %s, should be label2Value", tmpV)
			}
		}
	}
}

type MockNode struct {
	id         int64
	mockLabels []string
	mockProps  map[string]interface{}
}

func (m MockNode) Id() int64 {
	return m.id
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
