package internal

import (
	"errors"
	"reflect"
)

type Mapping struct {
	Label      string
	Attributes map[string]string
}

func GetMapping(input interface{}) (Mapping, error) {

	var mapping Mapping

	label := reflect.TypeOf(input)

	if label.Kind() != reflect.Struct {
		return mapping, errors.New("input is not a struct")
	}
	mapping.Label = label.Name()
	mapping.Attributes = make(map[string]string, label.NumField())

	for i := 0; i < label.NumField(); i++ {
		f := label.Field(i)
		mapping.Attributes[f.Name] = f.Type.Name()
	}

	return mapping, nil
}
