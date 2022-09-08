package utils

import (
	"encoding/json"
	"golang.org/x/exp/slices"
)

func Struct2ExpectedFieldMap(myStruct any, expectedFields []string) map[string]any {
	allFields := Struct2FieldMap(myStruct)

	expectedFieldsMap := make(map[string]any, len(allFields))
	for k, v := range allFields {
		if slices.Contains(expectedFields, k) {
			expectedFieldsMap[k] = v
		}
	}
	return expectedFieldsMap
}

func Struct2FieldMap(myStruct any) map[string]any {
	var inInterface map[string]any
	inrec, _ := json.Marshal(myStruct)
	json.Unmarshal(inrec, &inInterface)

	return inInterface
}
