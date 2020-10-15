package internal

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

// Mapping ...
type Mapping struct {
	Label      string
	Attributes map[string]string
	Values     map[string]interface{}
}

// Convert neo4j record to a struct
func Convert(targetStruct interface{}, props map[string]interface{}, attributes map[string]string) (interface{}, error) {
	rValue := reflect.ValueOf(targetStruct)
	fmt.Println("input Kind: ", rValue.Kind())
	fmt.Println("input Type: ", rValue.Type())
	if rValue.Kind() != reflect.Struct {
		err := errors.New("need a struct")
		return "", err
	}

	//make a copy
	copyValuePtr := reflect.New(rValue.Type())
	fmt.Println("copy Kind: ", copyValuePtr.Kind()) // => ptr
	fmt.Println("copy Type: ", copyValuePtr.Type()) //*(struct)

	// deference
	copyValue := copyValuePtr.Elem()
	fmt.Println("deferenced Kind: ", copyValue.Kind())
	fmt.Println("deferenced Type: ", copyValue.Type())

	nbField := copyValue.NumField()
	fmt.Println("nb fields : ", nbField)

	empty := reflect.Value{}
	for k, v := range props {
		fieldVal := copyValue.FieldByName(k)
		if fieldVal != empty {
			err := setValue(&fieldVal, attributes[k], v.(string))
			if err != nil {
				return "", err
			}
		}
	}
	return copyValue.Interface(), nil
}

// GetMapping builds a Mapping structure for a given struct
func GetMapping(input interface{}) (Mapping, error) {

	var mapping Mapping

	val := reflect.ValueOf(input)

	if val.Kind() != reflect.Struct {
		return mapping, errors.New("input is not a struct")
	}

	structType := val.Type()
	mapping.Label = structType.Name()
	mapping.Attributes = make(map[string]string, structType.NumField())
	mapping.Values = make(map[string]interface{}, structType.NumField())

	for i := 0; i < val.NumField(); i++ {
		fieldType := val.Type().Field(i)
		if fieldType.Name != "Model" {
			mapping.Attributes[fieldType.Name] = fieldType.Type.Name()

			val, err := getValue(val.Field(i))
			if err != nil {
				return mapping, err
			}
			mapping.Values[fieldType.Name] = val
		}
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
func setValue(v *reflect.Value, targetType string, value string) error {
	switch targetType {
	case "string":
		v.SetString(value)
		return nil
	case "int":
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		v.SetInt(int64(intVal))
		return nil
	case "float64":
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		v.SetFloat(floatVal)
		return nil
	default:
		msg := fmt.Sprintf("setValue for %v is not implemented", targetType)
		return errors.New(msg)
	}
}
