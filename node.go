package whiterabbit

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/neo4j/neo4j-go-driver/neo4j"
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
	expectedType := node.Labels()[0]
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
	for k, v := range node.Props() {
		// TODO: decide if we keep this matching between neo4j and my struct
		upperK := strings.ToUpper(k[:1]) + k[1:]
		fieldVal := copyValue.FieldByName(upperK)
		if (fieldVal == reflect.Value{}) {
			addLabels[k] = v.(string)
		} else {
			fieldVal := copyValue.FieldByName(upperK)
			err := setValue2(&fieldVal, v.(string))
			if err != nil {
				return "", err
			}
		}
	}
	model := copyValue.FieldByName("Model")
	if (model != reflect.Value{}) {
		idField := model.FieldByName("ID")
		idField.SetInt(node.Id())
		labField := model.FieldByName("Labels")
		addLabelsValue := reflect.ValueOf(addLabels)
		labField.Set(addLabelsValue)
	}
	return copyValue.Interface(), nil
}
func setValue2(fv *reflect.Value, value string) error {
	switch fv.Kind() {
	case reflect.String:
		fv.SetString(value)
	case reflect.Int:
		intVal, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		fv.SetInt(int64(intVal))
	case reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		fv.SetFloat(floatVal)
		return nil
	default:
		msg := fmt.Sprintf("setValue for %s is not implemented", fv.Kind())
		return errors.New(msg)
	}
	return nil
}
