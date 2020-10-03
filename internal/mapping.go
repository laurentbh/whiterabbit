package internal

import (
	"errors"
	"reflect"
)

type Mapping struct {
	Label      string
	Attributes map[string]string
	Values     map[string]interface{}
}

func GetMapping(input interface{}) (Mapping, error) {

	var mapping Mapping

	valType := reflect.ValueOf(input)

	if valType.Kind() != reflect.Struct {
		return mapping, errors.New("input is not a struct")
	}
	val := reflect.ValueOf(input).Elem()

	structType := valType.Type()
	mapping.Label = structType.Name()
	mapping.Attributes = make(map[string]string, structType.NumField())
	mapping.Values = make(map[string]interface{}, structType.NumField())

	for i := 0; i < valType.NumField(); i++ {
		fieldType := valType.Type().Field(i)
		fieldValue := val.Field(i)
		mapping.Attributes[fieldType.Name] = fieldType.Type.Name()
		mapping.Values[fieldType.Name] = fieldValue
	}

	return mapping, nil
}
