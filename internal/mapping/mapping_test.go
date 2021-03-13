package mapping

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type test struct {
	A int
	B string
	C float64
}

func TestMappingWithEmptyField(t *testing.T) {

	input := test{B: "i am a string"}
	expected := Mapping{
		Label:      "Anode",
		Attributes: map[string]string{"A": "int", "B": "string", "C": "float64"},
		Values:     map[string]interface{}{"B": "i am a string"},
	}
	mapping, err := GetMapping(input)
	assert.Nil(t, err, "should not return error")

	for k := range expected.Attributes {
		_, ok := mapping.Attributes[k]
		if !ok {
			t.Errorf("missing attribute [%s]", k)
		} else {
			assert.Equal(t, expected.Attributes[k], mapping.Attributes[k])
		}
	}
}
func TestMapping(t *testing.T) {

	expected := Mapping{
		Label:      "test",
		Attributes: map[string]string{"A": "int", "B": "string", "C": "float64"},
		Values:     map[string]interface{}{"A": "415", "B": "i am a string", "C": "3.14"},
	}

	input := test{A: 415, B: "i am a string", C: 3.14}
	mapping, err := GetMapping(input)

	assert.Nil(t, err, "should not return error")

	assert.Equal(t, expected.Label, mapping.Label, "wrong label")

	for k := range expected.Attributes {
		_, ok := mapping.Attributes[k]
		if !ok {
			t.Errorf("missing attribute [%s]", k)
		} else {
			assert.Equal(t, expected.Attributes[k], mapping.Attributes[k])
		}

	}
}
func TestMappingError(t *testing.T) {
	var i int
	_, err := GetMapping(i)
	assert.Nil(t, err)
}
