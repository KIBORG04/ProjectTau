package utils

import (
	"math/rand"
	"reflect"
)

func IsSlice(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Slice
}

func Slice2float32Map(slice []string) map[string]float32 {
	elementMap := make(map[string]float32)
	for i := 0; i < len(slice); i++ {
		elementMap[slice[i]] = 0
	}
	return elementMap
}

func Slice2int32Map(slice []string) map[string]int32 {
	elementMap := make(map[string]int32)
	for i := 0; i < len(slice); i++ {
		elementMap[slice[i]] = 0
	}
	return elementMap
}

func Contains(s interface{}, str interface{}) bool {
	if IsSlice(s) {
		v := reflect.ValueOf(s)
		for i := 0; i < v.Len(); i++ {
			if v.Index(i).Interface() == str {
				return true
			}
		}
	}
	return false
}

func Pick(from interface{}) interface{} {
	if IsSlice(from) {
		v := reflect.ValueOf(from)
		return v.Index(rand.Intn(v.Len())).Interface()
	}
	return nil
}
