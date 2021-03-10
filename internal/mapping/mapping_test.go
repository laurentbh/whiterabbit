package mapping

import "testing"

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
	if err != nil {
		t.Error("GetMqapping error ", err)
	}
	for k := range expected.Attributes {
		_, ok := mapping.Attributes[k]
		if !ok {
			t.Errorf("missing attribute [%s]", k)
		} else if mapping.Attributes[k] != expected.Attributes[k] {
			t.Errorf("attribute error key:[%s] expecting [%s] got [%s]", k, expected.Attributes[k], mapping.Attributes[k])
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
	if err != nil {
		t.Error("GetMqapping error ", err)
	}

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
