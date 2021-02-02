package whiterabbit

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// ConvertNode converts neo4j node into one of the candidate struct
func ConvertNode(node neo4j.Node, candidates []interface{}) (interface{}, error) {

	var candidateType []string
	// build list of type while making sure all candidates are struct
	for _, c := range candidates {
		rValue := reflect.ValueOf(c)
		if rValue.Kind() != reflect.Struct {
			return "", fmt.Errorf("convertNode: %v is not a struct", c)
		}
		tmpType := rValue.Type().Name()
		if idx := strings.LastIndex(tmpType, "."); idx != -1 {
			candidateType = append(candidateType, tmpType[idx+1:])

		} else {
			candidateType = append(candidateType, tmpType)
		}
	}
	// verify node is in the candidates
	// TODO: handle nodes with multiple labels
	expectedType := node.Labels[0]
	var candIdx int = -1

	for i, c := range candidateType {
		if c == expectedType {
			candIdx = i
			break
		}
	}
	if candIdx == -1 {
		return "", fmt.Errorf("node %s is not list of candidates", expectedType)
	}

	rValue := reflect.ValueOf(candidates[candIdx])
	//make a copy
	copyValuePtr := reflect.New(rValue.Type())
	copyValue := copyValuePtr.Elem()

	// prepare additional labels (anything not in rValue)
	addLabels := make(map[string]string)
	for k, v := range node.Props {
		// TODO: decide if we keep this matching between neo4j and my struct
		upperK := strings.ToUpper(k[:1]) + k[1:]
		fieldVal := copyValue.FieldByName(upperK)
		if (fieldVal == reflect.Value{}) {
			addLabels[k] = v.(string)
		} else {
			fieldVal := copyValue.FieldByName(upperK)
			err := setValueNeoToStruct(&fieldVal, v)
			if err != nil {
				return "", err
			}
		}
	}
	model := copyValue.FieldByName("Model")
	if (model != reflect.Value{}) {
		idField := model.FieldByName("ID")
		idField.SetInt(node.Id)
		labField := model.FieldByName("Labels")
		addLabelsValue := reflect.ValueOf(addLabels)
		labField.Set(addLabelsValue)
	}
	return copyValue.Interface(), nil
}
func setValueNeoToStruct(fv *reflect.Value, value interface{}) error {
	switch fv.Kind() {
	case reflect.String:
		fv.SetString(value.(string))
	case reflect.Int:
		// neo always encodes int as int64
		conv, ok := value.(int64)
		if ok {
			fv.SetInt(int64(conv))
		} else {
			return fmt.Errorf("can't convert [%v] to int", value)
		}
		return nil
	case reflect.Float64:
		conv, ok := value.(float64)
		if ok {
			fv.SetFloat(conv)
		} else {
			// edge case, if float value as no decimal part, it's returned as int by neo4j
			if reflect.TypeOf(value).Kind() == reflect.Int64 {
				conv2, ok := value.(int64)
				if ok {
					fv.SetFloat(float64(conv2))
					return nil
				}

			}
			return fmt.Errorf("can't convert [%v] to float", value)
		}
		return nil
	case reflect.Int64:
		conv, ok := value.(int64)
		if ok {
			fv.SetInt(conv)
		} else {
			return fmt.Errorf("can't convert [%v] to int64", value)
		}
		return nil
	case reflect.Slice:
		v2 := reflect.ValueOf(value)
		// length of value
		len := v2.Len()
		if len == 0 {
			return nil
		}
		first := v2.Index(0)
		// type of inner type
		switch reflect.TypeOf(first.Interface()).Kind() {
		case reflect.String:
			// alloc slice of strings
			fv.Set( reflect.MakeSlice(reflect.TypeOf([]string{}), len, len))
			// copy elements
			for i:=0; i<len; i++ {
				fv.Index(i).Set(reflect.ValueOf(v2.Index(i).Interface()))
			}
			return nil
		default:
			msg := fmt.Sprintf("setValue for %s of %s is not implemented", fv.Kind(),
				reflect.TypeOf(first.Interface()).Kind())
			return errors.New(msg)
		}
	default:
		msg := fmt.Sprintf("setValue for %s is not implemented", fv.Kind())
		return errors.New(msg)
	}
	return nil
}
