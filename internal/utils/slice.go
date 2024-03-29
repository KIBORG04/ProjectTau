package utils

import (
	"golang.org/x/exp/slices"
	"math/rand"

	"golang.org/x/exp/constraints"
)

func Slice2Map[T constraints.Integer](slice []string) map[string]T {
	elementMap := make(map[string]T)
	for i := 0; i < len(slice); i++ {
		elementMap[slice[i]] = 0
	}
	return elementMap
}

func Pick[T any](from []T) T {
	return from[rand.Intn(len(from))]
}

func RemoveElem[s comparable](slice []s, elem s) []s {
	indx := slices.Index(slice, elem)
	slice[indx] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}

func RemoveByIndex[T comparable](slice []T, index int) []T {
	slice[index] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}

func GetKeyByValue[T, R comparable](myMap map[R]T, el T) R {
	var genericError R
	for k, v := range myMap {
		if el == v {
			return k
		}
	}
	return genericError
}
