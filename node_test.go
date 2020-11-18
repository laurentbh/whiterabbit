package whiterabbit

import (
	"encoding/json"
	"testing"
)

func TestJSON(t *testing.T) {

	type testForJSON struct {
		Name  string `json:"name"`
		Age   int64  `json:"age"`
		Model `json:"model"`
	}

	m := Model{ID: 475, Labels: map[string]string{"one": "1", "two": "2"}}

	test := testForJSON{Model: m, Name: "nobody", Age: 12}

	bytePayload, err := json.Marshal(test)
	if err != nil {
		t.Errorf("shouldn't produce error [%s]", err)
	}
	var res testForJSON
	err = json.Unmarshal(bytePayload, &res)

	if res.Model.ID != 475 || res.Model.Labels["one"] != "1" {
		t.Errorf("error processing")
	}
}
func TestConvertNode(t *testing.T) {

	mock := MockNode{
		id:         6,
		mockLabels: []string{"TestStruct"},
		mockProps: map[string]interface{}{
			"A":        (int64)(123), // neo4j returns all int as int64
			"B":        "valueForB",
			"C":        3.14,
			"newLabel": "newValue",
			"label2":   "label2Value"},
	}
	var candidate []interface{}
	candidate = append(candidate, TestStruct{})
	ret, err := ConvertNode(mock, candidate)
	if err != nil {
		t.Errorf("TestConvertNode : %s", err)
	} else {
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
	return m.mockProps
}
