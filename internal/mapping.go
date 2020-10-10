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
func Convert(targetStruct interface{}, props map[string]interface{}, attributes map[string]string) interface{} {
	rValue := reflect.ValueOf(targetStruct)
	fmt.Println("input Kind: ", rValue.Kind())
	fmt.Println("input Type: ", rValue.Type())
	// TODO: make sure it's a struct and not a pointer

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
			setValue(&fieldVal, attributes[k], v)
		}
	}
	return copyValue.Interface()
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
func setValue(v *reflect.Value, targetType string, value interface{}) error {
	switch targetType {
	case "string":
		v.SetString(value.(string))
		return nil
	case "int":
		fmt.Printf("int kind = %s", v.Kind().String())
		v.SetInt(int64(value.(int)))
		return nil
	case "float64":
		v.SetFloat(value.(float64))
		return nil
	default:
		msg := fmt.Sprintf("setValue for %v is not implemented", targetType)
		return errors.New(msg)
	}
}
