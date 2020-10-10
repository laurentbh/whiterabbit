package internal

import (
	"fmt"
	"reflect"
	"testing"
)

type TestStruct struct {
	A int
	B string
	C float64
}

func TestLol2(t *testing.T) {

	props := map[string]string{"A": "123", "B": "multi string", "C": "3.14"}
	attrs := map[string]string{"A": "int", "B": "string", "C": "float64"}

	ret := lol2(TestStruct{}, props, attrs)
	kind := reflect.TypeOf(ret)
	fmt.Printf("kind =%v", kind)

	c := ret.(TestStruct)
	if c.A != 123 || c.B != "multi string" || c.C != 3.142 {
		t.Errorf("error ret is %#v", c)
	}
}

func TestMappingWithEmptyField(t *testing.T) {

	input := TestStruct{B: "i am a string"}
	expected := Mapping{
		Label:      "Anode",
		Attributes: map[string]string{"A": "int", "B": "string", "C": "float64"},
		Values:     map[string]interface{}{"B": "i am a string"},
	}
	mapping, _ := GetMapping(input)
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
		Label:      "Anode",
		Attributes: map[string]string{"A": "int", "B": "string", "C": "float64"},
		Values:     map[string]interface{}{"A": "415", "B": "i am a string", "C": "3.14"},
	}

	input := TestStruct{A: 415, B: "i am a string", C: 3.14}
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
