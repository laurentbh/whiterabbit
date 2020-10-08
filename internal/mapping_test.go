package internal

import (
	"strings"
	"testing"
)

type test struct {
	name string
	id   uint64
}

func TestMapping(t *testing.T) {

	type Anode struct {
		A int
		B string
		C float64
	}

	expected := Mapping{
		Label:      "Anode",
		Attributes: map[string]string{"A": "int", "B": "string", "C": "float64"},
		Values:     map[string]string{"A": "415", "B": "i am a sstring", "C": "3.14"},
	}

	input := Anode{A: 415, B: "i am a string", C: 3.14}
	mapping, _ := GetMapping(input)

	// TODO : find sonething for assertions
	if mapping.Label != expected.Label {
		t.Errorf("label error, expecting [%s] got [%s]", expected.Label, mapping.Label)
	}

	for k := range expected.Attributes {
		_, ok := mapping.Attributes[k]
		if !ok {
			t.Errorf("missing attribute [%s]", k)
		} else if mapping.Attributes[k] != expected.Attributes[k] {
			t.Errorf("attribute error key:[%s] expecting [%s] got [%s]", k, expected.Attributes[k], mapping.Attributes[k])
		} else if !strings.HasPrefix(mapping.Values[k], expected.Values[k]) {
			t.Errorf("value error key:[%s] expecting [%s] got [%s]", k, expected.Values[k], mapping.Values[k])

		}

	}
}
func TestMappingError(t *testing.T) {
	var i int
	_, err := GetMapping(i)
	if err == nil {
		t.Error("expecting")
	}

}
