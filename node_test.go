package whiterabbit

import (
	"encoding/json"
	"testing"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/stretchr/testify/assert"
)

type test struct {
	Model
	A int
	B string
	C float64
}

type test2 struct {
	Model
	A int
}

func TestJSON(t *testing.T) {

	type testForJSON struct {
		Name  string `json:"name"`
		Age   int64  `json:"age"`
		Model `json:"model"`
	}

	m := Model{ID: 475, Labels: map[string]string{"one": "1", "two": "2"}}

	test := testForJSON{Model: m, Name: "nobody", Age: 12}

	bytePayload, err := json.Marshal(test)
	assert.Nil(t, err, "should not error")

	var res testForJSON
	err = json.Unmarshal(bytePayload, &res)

	var expectedID int64 = 475
	assert.Equal(t, expectedID, res.Model.ID, "wrong value for Model ID")
	assert.Equal(t, "1", res.Model.Labels["one"], "wrong value for Model Label")
}
func TestConvertNode(t *testing.T) {
	mock := neo4j.Node{
		Id:     6,
		Labels: []string{"test"},
		Props: map[string]interface{}{
			"A":        (int64)(123), // neo4j returns all int as int64
			"B":        "valueForB",
			"C":        3.14,
			"newLabel": "newValue",
			"label2":   "label2Value"},
	}
	ret, err := ConvertNode(mock, test{})
	assert.Nil(t, err, "should not return error")
	assert.IsType(t, test{}, ret, "wrong return type")

	assert.Equal(t, 123, (ret.(test)).A, "wrong value for filed A")
	assert.Equal(t, "valueForB", (ret.(test)).B, "wrong value for field B")
	assert.Equal(t, 3.14, (ret.(test)).C, "wrong value for field C")

	var expectedID int64 = 6
	assert.Equal(t, expectedID, (ret.(test)).Model.ID, "wrong value for ID")
	assert.Equal(t, "newValue", (ret.(test)).Labels["newLabel"], "wrong value for newLabel")
	assert.Equal(t, "label2Value", (ret.(test)).Labels["label2"], "wrong value for label2Value")
}
func TestConvertNodeError(t *testing.T) {
	mock := neo4j.Node{
		Id:     6,
		Labels: []string{"test"},
		Props: map[string]interface{}{
			"A":        (int64)(123), // neo4j returns all int as int64
			"B":        "valueForB",
			"C":        3.14,
			"newLabel": "newValue",
			"label2":   "label2Value"},
	}
	_, err := ConvertNode(mock, 1, "should return error")
	assert.NotNil(t, err)

	_, err = ConvertNode(mock, test{}, 1, "should return error")
	assert.NotNil(t, err)
}
func TestConvertNodeMultipleCandidate(t *testing.T) {
	mock := neo4j.Node{
		Id:     6,
		Labels: []string{"test"},
		Props: map[string]interface{}{
			"A":        (int64)(123), // neo4j returns all int as int64
			"B":        "valueForB",
			"C":        3.14,
			"newLabel": "newValue",
			"label2":   "label2Value"},
	}
	ret, err := ConvertNode(mock, test2{}, test{})

	assert.Nil(t, err, "should not return error")
	assert.IsType(t, test{}, ret, "wrong return type")
}
func TestConvertFloat(t *testing.T) {
	var val int64
	val = 3
	mock := neo4j.Node{
		Id:     6,
		Labels: []string{"test"},
		Props: map[string]interface{}{
			"C": val},
	}
	ret, err := ConvertNode(mock, test{})
	assert.Nil(t, err, "should not return error")
	assert.IsType(t, test{}, ret, "wrong return type")
	var expectedFloat float64 = 3
	assert.Equal(t, expectedFloat, (ret.(test)).C, "wrong value for field C")
}
