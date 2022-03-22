package utils

import (
	"golang.org/x/exp/slices"
	"reflect"
	"strings"
)

func JsonFieldNames(v interface{}, expectedFields *[]string) []string {
	names := make([]string, 0)
	typeof := reflect.TypeOf(v).Elem()
	for i := 0; i < typeof.NumField(); i++ {
		field := typeof.Field(i)
		fieldName := JsonFieldName(field)
		if expectedFields != nil && !slices.Contains(*expectedFields, fieldName) {
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
