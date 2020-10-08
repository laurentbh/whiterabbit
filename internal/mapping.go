package internal

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type Mapping struct {
	Label      string
	Attributes map[string]string
	Values     map[string]string
}

func GetMapping(input interface{}) (Mapping, error) {

	var mapping Mapping

	val := reflect.ValueOf(input)

	if val.Kind() != reflect.Struct {
		return mapping, errors.New("input is not a struct")
	}

	structType := val.Type()
	mapping.Label = structType.Name()
	mapping.Attributes = make(map[string]string, structType.NumField())
	mapping.Values = make(map[string]string, structType.NumField())

	for i := 0; i < val.NumField(); i++ {
		fieldType := val.Type().Field(i)
		mapping.Attributes[fieldType.Name] = fieldType.Type.Name()

		val, err := getValue(val.Field(i))
		if err != nil {
			return mapping, err
		}
		mapping.Values[fieldType.Name] = val
	}

	return mapping, nil
}

func getValue(v reflect.Value) (string, error) {
	switch v.Kind() {
	case reflect.Int:
		return strconv.Itoa(v.Interface().(int)), nil
	case reflect.String:
		return v.Interface().(string), nil
	case reflect.Float64:
		fl := v.Interface().(float64)
		return fmt.Sprintf("%f", fl), nil
	}
	msg := fmt.Sprintf("getValue for %v is not implemented", v.Kind())
	return "", errors.New(msg)
}
