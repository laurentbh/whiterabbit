package internal

import (
	"reflect"
	"testing"
)

type TestStruct struct {
	Model
	A int
	B string
	C float64
}

func TestConvert(t *testing.T) {

	props := map[string]interface{}{"A": "123", "B": "multi string", "C": "3.14"}
	attrs := map[string]string{"A": "int", "B": "string", "C": "float64"}

	ret, _ := Convert(TestStruct{}, props, attrs)
	retType := reflect.TypeOf(ret)

	if retType != reflect.TypeOf(TestStruct{}) {
		t.Errorf("returned type is %s, expecting %s", retType, reflect.TypeOf(TestStruct{}))
	}

	c := ret.(TestStruct)
	if c.A != 123 || c.B != "multi string" || c.C != 3.14 {
		t.Errorf("error ret is %#v", c)
	}
}
func TestConvertBadArg(t *testing.T) {
	props := map[string]interface{}{"A": "123", "B": "multi string", "C": "3.14"}
	attrs := map[string]string{"A": "int", "B": "string", "C": "float64"}

	_, err := Convert(&TestStruct{}, props, attrs)
	if err == nil {
		t.Errorf("expecting error")
	}
	_, err = Convert(10, props, attrs)
	if err == nil {
		t.Errorf("expecting error")
	}
}
func TestConvertUnHandled(t *testing.T) {
	type dummy struct {
		A int
		B []string
	}
	props := map[string]interface{}{"A": "123", "B": "[]string{\"str1\", \"str2\"}"}
	attrs := map[string]string{"A": "int", "B": "[]string"}

	_, err := Convert(dummy{}, props, attrs)
	if err == nil {
		t.Errorf("expecting error")
	}
}

func TestMappingWithEmptyField(t *testing.T) {

	input := TestStruct{B: "i am a string"}
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
		Label:      "TestStruct",
		Attributes: map[string]string{"A": "int", "B": "string", "C": "float64"},
		Values:     map[string]interface{}{"A": "415", "B": "i am a string", "C": "3.14"},
	}

	input := TestStruct{A: 415, B: "i am a string", C: 3.14}
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
