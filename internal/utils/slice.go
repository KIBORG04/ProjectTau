package utils

import (
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
