package utils

import (
	"encoding/json"

	"gorm.io/gorm/utils"
)

func Struct2ExpectedFieldMap(myStruct interface{}, expectedFields []string) map[string]interface{} {
	allFields := Struct2FieldMap(myStruct)

	expectedFieldsMap := make(map[string]interface{}, len(allFields))
	for k, v := range allFields {
		if utils.Contains(expectedFields, k) {
			expectedFieldsMap[k] = v
		}
	}
	return expectedFieldsMap
}

func Struct2FieldMap(myStruct interface{}) map[string]interface{} {
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(myStruct)
	json.Unmarshal(inrec, &inInterface)

	return inInterface
}
