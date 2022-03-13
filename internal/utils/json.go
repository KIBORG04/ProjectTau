package utils

import (
	"fmt"
	"reflect"
	"strings"
)

var (
	Ints = []string{
		"int",
		"int8",
		"int16",
		"int32",
		"int64",
		"uint",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
	}
)

func JsonFieldNames(v interface{}, expectedFields *[]string, expectedTypes *[]string) []string {
	names := make([]string, 0)
	typeof := reflect.TypeOf(v).Elem()
	for i := 0; i < typeof.NumField(); i++ {
		field := typeof.Field(i)
		if expectedTypes != nil && !Contains(*expectedTypes, field.Type.Name()) {
			continue
		}

		fieldName := JsonFieldName(field)
		if expectedFields != nil && !Contains(*expectedFields, fieldName) {
			continue
		}
		if len(fieldName) == 0 {
			continue
		}
		names = append(names, fieldName)
	}

	return names
}

func JsonFieldName(field reflect.StructField) string {
	tag := field.Tag.Get("json")
	jsonFieldName := strings.Split(tag, ",")[0]
	return jsonFieldName
}

func FieldNameByValue(v interface{}) map[string]string {
	m := make(map[string]string, 0)

	val := reflect.ValueOf(v).Elem()
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)

		val := reflect.ValueOf(valueField.Interface())
		m[JsonFieldName(typeField)] = fmt.Sprintf("%v", val)
	}

	return m
}
